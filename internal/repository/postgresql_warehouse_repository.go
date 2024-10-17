package repository

import (
	"context"
	"errors"
	"simcomm-monolith/internal/model"
	"simcomm-monolith/util"

	log "github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type WarehouseRepository interface {
	Create(ctx context.Context, warehouse *model.Warehouse) error
	Get(ctx context.Context, id int) (*model.Warehouse, error)
	GetAll(ctx context.Context) ([]model.Warehouse, error)
	Update(ctx context.Context, warehouse *model.Warehouse) error
	Delete(ctx context.Context, id int) error

	WarehouseStoredProductRepository
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

type WarehouseStoredProductRepository interface {
	WSPCreate(ctx context.Context, warehousestoredproduct *model.WarehouseStoredProduct) error
	WSPGet(ctx context.Context, id int) (*model.WarehouseStoredProduct, error)
	WSPGetAll(ctx context.Context) ([]model.WarehouseStoredProduct, error)
	WSPUpdate(ctx context.Context, warehousestoredproduct *model.WarehouseStoredProduct) error
	WSPDelete(ctx context.Context, id int) error

	WSPGetByShopProductID(ctx context.Context, shopProductID int, warehouseID int) (*model.WarehouseStoredProduct, error)
	WSPSubstractStock(ctx context.Context, warehousestoredproduct *model.WarehouseStoredProduct, subtrahend int) error
}

// Create inserts a new warehousestoredproduct into the database
func (r *postgresWarehouseRepository) WSPCreate(ctx context.Context, warehousestoredproduct *model.WarehouseStoredProduct) error {
	return r.db.WithContext(ctx).Create(warehousestoredproduct).Error
}

// Get retrieves a warehousestoredproduct by ID
func (r *postgresWarehouseRepository) WSPGet(ctx context.Context, id int) (*model.WarehouseStoredProduct, error) {
	var warehousestoredproduct model.WarehouseStoredProduct
	if err := r.db.WithContext(ctx).First(&warehousestoredproduct, id).Error; err != nil {
		return nil, err
	}
	return &warehousestoredproduct, nil
}

// GetAll retrieves all warehousestoredproducts from the database
func (r *postgresWarehouseRepository) WSPGetAll(ctx context.Context) ([]model.WarehouseStoredProduct, error) {
	var warehousestoredproducts []model.WarehouseStoredProduct
	if err := r.db.WithContext(ctx).Find(&warehousestoredproducts).Error; err != nil {
		return nil, err
	}
	return warehousestoredproducts, nil
}

// Update updates an existing warehousestoredproduct
func (r *postgresWarehouseRepository) WSPUpdate(ctx context.Context, warehousestoredproduct *model.WarehouseStoredProduct) error {
	if err := r.db.WithContext(ctx).Save(warehousestoredproduct).Error; err != nil {
		return err
	}
	return nil
}

func (r *postgresWarehouseRepository) WSPSubstractStock(ctx context.Context, wsp *model.WarehouseStoredProduct, subtrahend int) error {
	wsp.Stock = wsp.Stock - subtrahend
	wsp.UpdatedAt = util.TimeNow()

	if err := r.db.WithContext(ctx).
		Where("stock > ? and id = ? ", subtrahend, wsp.ID).
		Save(wsp).Error; err != nil {

		return err
	}
	return nil
}

// Delete removes a warehousestoredproduct from the database
func (r *postgresWarehouseRepository) WSPDelete(ctx context.Context, id int) error {
	if err := r.db.WithContext(ctx).Delete(&model.WarehouseStoredProduct{}, id).Error; err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (r *postgresWarehouseRepository) WSPGetByShopProductID(ctx context.Context, shopProductID int, warehouseID int) (*model.WarehouseStoredProduct, error) {
	var warehousestoredproduct model.WarehouseStoredProduct
	if err := r.db.WithContext(ctx).
		Where("warehouse_id = ? AND shop_product_id = ?", warehouseID, shopProductID).
		First(&warehousestoredproduct).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &warehousestoredproduct, nil
}
