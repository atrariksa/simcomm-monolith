package repository

import (
	"context"
	"simcomm-monolith/internal/model"

	log "github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type ShopRepository interface {
	Create(ctx context.Context, shop *model.Shop) error
	Get(ctx context.Context, id int) (*model.Shop, error)
	GetAll(ctx context.Context) ([]model.Shop, error)
	Update(ctx context.Context, shop *model.Shop) error
	Delete(ctx context.Context, id int) error
}

type postgresShopRepository struct {
	db *gorm.DB
}

// NewShopRepository creates a new instance of ShopRepository
func NewPostgreShopRepository(db *gorm.DB) *postgresShopRepository {
	return &postgresShopRepository{db: db}
}

// Create inserts a new shop into the database
func (r *postgresShopRepository) Create(ctx context.Context, shop *model.Shop) error {
	return r.db.WithContext(ctx).Create(shop).Error
}

// Get retrieves a shop by ID
func (r *postgresShopRepository) Get(ctx context.Context, id int) (*model.Shop, error) {
	var shop model.Shop
	if err := r.db.WithContext(ctx).First(&shop, id).Error; err != nil {
		return nil, err
	}
	return &shop, nil
}

// GetAll retrieves all shops from the database
func (r *postgresShopRepository) GetAll(ctx context.Context) ([]model.Shop, error) {
	var shops []model.Shop
	if err := r.db.WithContext(ctx).Find(&shops).Error; err != nil {
		return nil, err
	}
	return shops, nil
}

// Update updates an existing shop
func (r *postgresShopRepository) Update(ctx context.Context, shop *model.Shop) error {
	if err := r.db.WithContext(ctx).Save(shop).Error; err != nil {
		return err
	}
	return nil
}

// Delete removes a shop from the database
func (r *postgresShopRepository) Delete(ctx context.Context, id int) error {
	if err := r.db.WithContext(ctx).Delete(&model.Shop{}, id).Error; err != nil {
		log.Error(err)
		return err
	}
	return nil
}
