package repository

import (
	"context"
	"simcomm-monolith/internal/model"

	log "github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) error
	Get(ctx context.Context, id int) (*model.Order, error)
	GetAll(ctx context.Context) ([]model.Order, error)
	Update(ctx context.Context, order *model.Order) error
	Delete(ctx context.Context, id int) error
}

type postgresOrderRepository struct {
	db *gorm.DB
}

// NewOrderRepository creates a new instance of OrderRepository
func NewPostgreOrderRepository(db *gorm.DB) *postgresOrderRepository {
	return &postgresOrderRepository{db: db}
}

// Create inserts a new order into the database
func (r *postgresOrderRepository) Create(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

// Get retrieves a order by ID
func (r *postgresOrderRepository) Get(ctx context.Context, id int) (*model.Order, error) {
	var order model.Order
	if err := r.db.WithContext(ctx).First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// GetAll retrieves all orders from the database
func (r *postgresOrderRepository) GetAll(ctx context.Context) ([]model.Order, error) {
	var orders []model.Order
	if err := r.db.WithContext(ctx).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

// Update updates an existing order
func (r *postgresOrderRepository) Update(ctx context.Context, order *model.Order) error {
	if err := r.db.WithContext(ctx).Save(order).Error; err != nil {
		return err
	}
	return nil
}

// Delete removes a order from the database
func (r *postgresOrderRepository) Delete(ctx context.Context, id int) error {
	if err := r.db.WithContext(ctx).Delete(&model.Order{}, id).Error; err != nil {
		log.Error(err)
		return err
	}
	return nil
}
