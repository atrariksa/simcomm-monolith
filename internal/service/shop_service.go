package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"simcomm-monolith/config"
	"simcomm-monolith/internal/model"
	"simcomm-monolith/internal/repository"
	"simcomm-monolith/util"

	log "github.com/labstack/gommon/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

// ShopService defines the methods for the Shop service
type ShopService interface {
	Create(ctx context.Context, shop *model.Shop) error
	Get(ctx context.Context, id int) (*model.Shop, error)
	GetAll(ctx context.Context) ([]model.Shop, error)
	Update(ctx context.Context, shop *model.Shop) error
	Delete(ctx context.Context, id int) error

	ShopProductService

	CreateTransferProduct(ctx context.Context, tp *model.TransferProduct) error
	ProcessRTPQueue(ctx context.Context, msg amqp.Delivery) error
	UpdateTransferProduct(ctx context.Context, utp *model.UpdateTransferProduct) error
}

type shopService struct {
	wspSvc    WarehouseService
	repo      repository.ShopRepository
	redisRepo repository.RedisRepository
	queues    map[string]repository.Queue
	cfg       *config.Config
}

func NewShopService(wspSvc WarehouseService, repo repository.ShopRepository, redisRepo repository.RedisRepository, q map[string]repository.Queue, cfg *config.Config) *shopService {
	return &shopService{
		wspSvc:    wspSvc,
		repo:      repo,
		redisRepo: redisRepo,
		queues:    q,
		cfg:       cfg,
	}
}

func (s *shopService) Create(ctx context.Context, shop *model.Shop) error {
	timeNow := util.TimeNow()
	shop.CreatedAt = timeNow
	shop.UpdatedAt = timeNow
	err := s.repo.Create(ctx, shop)
	if err != nil {
		log.Error(err)
	}
	return err
}

func (s *shopService) Get(ctx context.Context, id int) (*model.Shop, error) {
	return s.repo.Get(ctx, id)
}

func (s *shopService) GetAll(ctx context.Context) ([]model.Shop, error) {
	return s.repo.GetAll(ctx)
}

func (s *shopService) Update(ctx context.Context, shop *model.Shop) error {
	return s.repo.Update(ctx, shop)
}

func (s *shopService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

// ShopProductService defines the methods for the ShopProduct service
type ShopProductService interface {
	ShopProductServiceCreate(ctx context.Context, shopproduct *model.ShopProduct) error
	ShopProductServiceGet(ctx context.Context, id int) (*model.ShopProduct, error)
	ShopProductServiceGetAll(ctx context.Context) ([]model.ShopProduct, error)
	ShopProductServiceUpdate(ctx context.Context, shopproduct *model.ShopProduct) error
	ShopProductServiceDelete(ctx context.Context, id int) error
}

func (s *shopService) ShopProductServiceCreate(ctx context.Context, shopproduct *model.ShopProduct) error {
	timeNow := util.TimeNow()
	shopproduct.CreatedAt = timeNow
	shopproduct.UpdatedAt = timeNow
	err := s.repo.ShopProductRepositoryCreate(ctx, shopproduct)
	if err != nil {
		log.Error(err)
	}
	return err
}

func (s *shopService) ShopProductServiceGet(ctx context.Context, id int) (*model.ShopProduct, error) {
	return s.repo.ShopProductRepositoryGet(ctx, id)
}

func (s *shopService) ShopProductServiceGetAll(ctx context.Context) ([]model.ShopProduct, error) {
	return s.repo.ShopProductRepositoryGetAll(ctx)
}

func (s *shopService) ShopProductServiceUpdate(ctx context.Context, shopproduct *model.ShopProduct) error {
	return s.repo.ShopProductRepositoryUpdate(ctx, shopproduct)
}

func (s *shopService) ShopProductServiceDelete(ctx context.Context, id int) error {
	return s.repo.ShopProductRepositoryDelete(ctx, id)
}

func (s *shopService) CreateTransferProduct(ctx context.Context, tp *model.TransferProduct) error {
	timeNow := util.TimeNow()
	chShopProduct := make(chan model.ShopProduct, 1)
	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		shopProduct, err := s.ShopProductServiceGet(egCtx, tp.ShopProductID)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		chShopProduct <- *shopProduct
		return nil
	})
	var shopProduct = <-chShopProduct

	chWSPSource := make(chan model.WarehouseStoredProduct, 1)
	eg.Go(func() error {
		wspSource, err := s.wspSvc.WSPGetByShopProductID(egCtx, tp.ShopProductID, tp.WarehouseIDSource)
		if err != nil {
			return err
		}
		chWSPSource <- *wspSource
		return nil
	})
	var wspSource = <-chWSPSource

	if err := eg.Wait(); err != nil {
		log.Error(err)
		return err
	}

	if shopProduct.ID < 1 {
		return errors.New("shop product not found")
	}

	if wspSource.ID < 1 || wspSource.Stock < tp.StockToTransfer {
		return errors.New("stored product not enough")
	}

	shopProduct.Stock = shopProduct.Stock - tp.StockToTransfer
	shopProduct.UpdatedAt = timeNow
	spDetails := shopProduct.Detail.ShopProductDetails
	for i, v := range spDetails {
		if v.WarehouseID == wspSource.ID {
			spDetails[i].Stock = v.Stock - tp.StockToTransfer
			break
		}
	}
	shopProduct.Detail.ShopProductDetails = spDetails

	tp.Status = "OTW"
	tp.Detail = model.TransferProductDetail{
		Histories: []model.TransferProductHostory{
			{
				Status:    tp.Status,
				Timestamp: timeNow,
			},
		},
	}

	err := s.repo.ShopProductRepositoryCreateTransferProduct(ctx, tp, &shopProduct, s.queues[repository.TPQueueName])
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (s *shopService) ProcessRTPQueue(ctx context.Context, msg amqp.Delivery) error {
	var rtp model.UpdateTransferProduct
	err := json.Unmarshal(msg.Body, &rtp)
	if err != nil {
		log.Error(err)
		return err
	}

	eg, egCtx := errgroup.WithContext(ctx)

	chShopProduct := make(chan model.ShopProduct, 1)
	eg.Go(func() error {
		shopProduct, err := s.ShopProductServiceGet(egCtx, rtp.ShopProductID)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		chShopProduct <- *shopProduct
		return nil
	})
	var shopProduct = <-chShopProduct

	chTP := make(chan model.TransferProduct, 1)
	eg.Go(func() error {
		tp, err := s.repo.ShopProductRepositoryGetTransferProduct(egCtx, rtp.TransferProductID)
		if err != nil {
			return err
		}
		chTP <- *tp
		return nil
	})
	var tp = <-chTP

	if err := eg.Wait(); err != nil {
		log.Error(err)
		return err
	}

	if shopProduct.ID < 1 {
		log.Error(errors.New("shop product not found"))
		return nil
	}

	if tp.ID < 1 {
		log.Error(errors.New("data transfer product not found"))
		return nil
	}

	timeNow := util.TimeNow()

	tp.UpdatedAt = timeNow
	tp.Status = model.TransferProductStatus.Failed
	tp.Detail.Histories = append(tp.Detail.Histories, model.TransferProductHostory{
		Timestamp: timeNow,
		Status:    tp.Status,
		Note:      rtp.Note,
	})

	shopProduct.Stock = shopProduct.Stock + tp.StockToTransfer
	shopProduct.UpdatedAt = timeNow

	err = s.repo.ShopProductRepositoryUpdateTransferProduct(ctx, &tp, &shopProduct, s.queues[repository.UTPQueueName])
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (s *shopService) UpdateTransferProduct(ctx context.Context, utp *model.UpdateTransferProduct) error {
	tp, err := s.repo.ShopProductRepositoryGetTransferProduct(ctx, utp.TransferProductID)
	if err != nil {
		return err
	}

	if tp.ID < 1 {
		log.Error(errors.New("data transfer product not found"))
		return nil
	}

	timeNow := util.TimeNow()
	tp.UpdatedAt = timeNow
	tp.Status = utp.Status
	tp.Detail.Histories = append(tp.Detail.Histories, model.TransferProductHostory{
		Timestamp: timeNow,
		Status:    tp.Status,
		Note:      utp.Note,
	})
	shopProduct, err := s.ShopProductServiceGet(ctx, utp.ShopProductID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if shopProduct.ID < 1 {
		err = fmt.Errorf("shop product not found : %v", utp.ShopProductID)
		log.Error(err)
		return err
	}

	if utp.Status == model.TransferProductStatus.Completed {
		spDetails := shopProduct.Detail.ShopProductDetails
		for i, v := range spDetails {
			if v.WarehouseID == tp.WarehouseIDDestination {
				spDetails[i].Stock = v.Stock + tp.StockToTransfer
				break
			}
		}
		shopProduct.Detail.ShopProductDetails = spDetails
		err = s.repo.ShopProductRepositoryUpdateTransferProduct(ctx, tp, shopProduct, s.queues[repository.UTPQueueName])
		if err != nil {
			log.Error(err)
			return err
		}
		return nil
	}

	spDetails := shopProduct.Detail.ShopProductDetails
	for i, v := range spDetails {
		if v.WarehouseID == tp.WarehouseIDSource {
			spDetails[i].Stock = v.Stock + tp.StockToTransfer
			break
		}
	}
	shopProduct.Detail.ShopProductDetails = spDetails

	if utp.Status == model.TransferProductStatus.Canceled {
		err = s.repo.ShopProductRepositoryUpdateTransferProduct(ctx, tp, shopProduct, s.queues[repository.UTPQueueName])
		if err != nil {
			log.Error(err)
			return err
		}
		return nil
	}
	if utp.Status == model.TransferProductStatus.Failed {
		err = s.repo.ShopProductRepositoryUpdateTransferProduct(ctx, tp, shopProduct, nil)
		if err != nil {
			log.Error(err)
			return err
		}
		return nil
	}

	return errors.New("Unprocessed Request")
}
