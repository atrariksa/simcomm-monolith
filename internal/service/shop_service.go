package service

import (
	"context"
	"errors"
	"simcomm-monolith/config"
	"simcomm-monolith/internal/model"
	"simcomm-monolith/internal/repository"
	"simcomm-monolith/util"

	log "github.com/labstack/gommon/log"
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
}

type shopService struct {
	wspSvc    WarehouseService
	repo      repository.ShopRepository
	redisRepo repository.RedisRepository
	queue     repository.Queue
	cfg       *config.Config
}

func NewShopService(wspSvc WarehouseService, repo repository.ShopRepository, redisRepo repository.RedisRepository, q repository.Queue, cfg *config.Config) *shopService {
	return &shopService{
		wspSvc:    wspSvc,
		repo:      repo,
		redisRepo: redisRepo,
		queue:     q,
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

	// chWSPSource := make(chan model.WarehouseStoredProduct, 1)
	// eg.Go(func() error {
	// 	wspSource, err := s.wspSvc.WSPGetByShopProductID(egCtx, tp.ShopProductID, tp.WarehouseIDSource)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	chWSPSource <- *wspSource
	// 	return nil
	// })
	// var wspSource = <-chWSPSource

	if err := eg.Wait(); err != nil {
		log.Error(err)
		return err
	}

	if shopProduct.ID < 1 {
		return errors.New("shop product not found")
	}

	// if wspSource.ID < 1 || wspSource.Stock < tp.StockToTransfer {
	// 	return errors.New("stored product not enough")
	// }

	shopProduct.Stock = shopProduct.Stock - tp.StockToTransfer
	shopProduct.UpdatedAt = timeNow

	tp.Status = "OTW"
	tp.Detail = model.TransferProductDetail{
		Histories: []model.TransferProductHostory{
			{
				Status:    tp.Status,
				Timestamp: timeNow,
			},
		},
	}

	err := s.repo.ShopProductRepositoryCreateTransferProduct(ctx, tp, &shopProduct)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
