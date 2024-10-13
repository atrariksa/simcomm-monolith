package repository

import (
	"context"
	"simcomm-monolith/internal/model"

	log "github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type WarehouseRepository interface {
	Create(ctx context.Context, warehouse *model.Warehouse) error
	Get(ctx context.Context, id int) (*model.Warehouse, error)
	GetAll(ctx context.Context) ([]model.Warehouse, error)
	Update(ctx context.Context, warehouse *model.Warehouse) error
	Delete(ctx context.Context, id int) error
}

type postgresWarehouseRepository struct {
	db *gorm.DB
}

// NewWarehouseRepository creates a new instance of WarehouseRepository
func NewPostgreWarehouseRepository(db *gorm.DB) *postgresWarehouseRepository {
	return &postgresWarehouseRepository{db: db}
}

// Create inserts a new warehouse into the database
func (r *postgresWarehouseRepository) Create(ctx context.Context, warehouse *model.Warehouse) error {
	return r.db.WithContext(ctx).Create(warehouse).Error
}

// Get retrieves a warehouse by ID
func (r *postgresWarehouseRepository) Get(ctx context.Context, id int) (*model.Warehouse, error) {
	var warehouse model.Warehouse
	if err := r.db.WithContext(ctx).First(&warehouse, id).Error; err != nil {
		return nil, err
	}
	return &warehouse, nil
}

// GetAll retrieves all warehouses from the database
func (r *postgresWarehouseRepository) GetAll(ctx context.Context) ([]model.Warehouse, error) {
	var warehouses []model.Warehouse
	if err := r.db.WithContext(ctx).Find(&warehouses).Error; err != nil {
		return nil, err
	}
	return warehouses, nil
}

// Update updates an existing warehouse
func (r *postgresWarehouseRepository) Update(ctx context.Context, warehouse *model.Warehouse) error {
	if err := r.db.WithContext(ctx).Save(warehouse).Error; err != nil {
		return err
	}
	return nil
}

// Delete removes a warehouse from the database
func (r *postgresWarehouseRepository) Delete(ctx context.Context, id int) error {
	if err := r.db.WithContext(ctx).Delete(&model.Warehouse{}, id).Error; err != nil {
		log.Error(err)
		return err
	}
	return nil
}
