package service

import (
	"context"
	"simcomm-monolith/config"
	"simcomm-monolith/internal/model"
	"simcomm-monolith/internal/repository"
	"simcomm-monolith/util"

	log "github.com/labstack/gommon/log"
)

// WarehouseService defines the methods for the Warehouse service
type WarehouseService interface {
	Create(ctx context.Context, warehouse *model.Warehouse) error
	Get(ctx context.Context, id int) (*model.Warehouse, error)
	GetAll(ctx context.Context) ([]model.Warehouse, error)
	Update(ctx context.Context, warehouse *model.Warehouse) error
	Delete(ctx context.Context, id int) error

	WarehouseStoredProductService
}

type warehouseService struct {
	repo      repository.WarehouseRepository
	redisRepo repository.RedisRepository
	cfg       *config.Config
}

func NewWarehouseService(repo repository.WarehouseRepository, redisRepo repository.RedisRepository, cfg *config.Config) *warehouseService {
	return &warehouseService{
		repo:      repo,
		redisRepo: redisRepo,
		cfg:       cfg,
	}
}

func (s *warehouseService) Create(ctx context.Context, warehouse *model.Warehouse) error {
	timeNow := util.TimeNow()
	warehouse.CreatedAt = timeNow
	warehouse.UpdatedAt = timeNow
	err := s.repo.Create(ctx, warehouse)
	if err != nil {
		log.Error(err)
	}
	return err
}

func (s *warehouseService) Get(ctx context.Context, id int) (*model.Warehouse, error) {
	return s.repo.Get(ctx, id)
}

func (s *warehouseService) GetAll(ctx context.Context) ([]model.Warehouse, error) {
	return s.repo.GetAll(ctx)
}

func (s *warehouseService) Update(ctx context.Context, warehouse *model.Warehouse) error {
	return s.repo.Update(ctx, warehouse)
}

func (s *warehouseService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

// WarehouseStoredProductService defines the methods for the WarehouseStoredProduct service
type WarehouseStoredProductService interface {
	WSPCreate(ctx context.Context, warehousestoredproduct *model.WarehouseStoredProduct) error
	WSPGet(ctx context.Context, id int) (*model.WarehouseStoredProduct, error)
	WSPGetAll(ctx context.Context) ([]model.WarehouseStoredProduct, error)
	WSPUpdate(ctx context.Context, warehousestoredproduct *model.WarehouseStoredProduct) error
	WSPDelete(ctx context.Context, id int) error
}

func (s *warehouseService) WSPCreate(ctx context.Context, warehousestoredproduct *model.WarehouseStoredProduct) error {
	timeNow := util.TimeNow()
	warehousestoredproduct.CreatedAt = timeNow
	warehousestoredproduct.UpdatedAt = timeNow
	err := s.repo.WSPCreate(ctx, warehousestoredproduct)
	if err != nil {
		log.Error(err)
	}
	return err
}

func (s *warehouseService) WSPGet(ctx context.Context, id int) (*model.WarehouseStoredProduct, error) {
	return s.repo.WSPGet(ctx, id)
}

func (s *warehouseService) WSPGetAll(ctx context.Context) ([]model.WarehouseStoredProduct, error) {
	return s.repo.WSPGetAll(ctx)
}

func (s *warehouseService) WSPUpdate(ctx context.Context, warehousestoredproduct *model.WarehouseStoredProduct) error {
	return s.repo.WSPUpdate(ctx, warehousestoredproduct)
}

func (s *warehouseService) WSPDelete(ctx context.Context, id int) error {
	return s.repo.WSPDelete(ctx, id)
}
