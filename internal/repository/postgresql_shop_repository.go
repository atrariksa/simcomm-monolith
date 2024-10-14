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

	ShopProductRepository
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

type ShopProductRepository interface {
	ShopProductRepositoryCreate(ctx context.Context, shopproduct *model.ShopProduct) error
	ShopProductRepositoryGet(ctx context.Context, id int) (*model.ShopProduct, error)
	ShopProductRepositoryGetAll(ctx context.Context) ([]model.ShopProduct, error)
	ShopProductRepositoryUpdate(ctx context.Context, shopproduct *model.ShopProduct) error
	ShopProductRepositoryDelete(ctx context.Context, id int) error
}

// Create inserts a new shopproduct into the database
func (r *postgresShopRepository) ShopProductRepositoryCreate(ctx context.Context, shopproduct *model.ShopProduct) error {
	return r.db.WithContext(ctx).Create(shopproduct).Error
}

// Get retrieves a shopproduct by ID
func (r *postgresShopRepository) ShopProductRepositoryGet(ctx context.Context, id int) (*model.ShopProduct, error) {
	var shopproduct model.ShopProduct
	if err := r.db.WithContext(ctx).First(&shopproduct, id).Error; err != nil {
		return nil, err
	}
	return &shopproduct, nil
}

// GetAll retrieves all shopproducts from the database
func (r *postgresShopRepository) ShopProductRepositoryGetAll(ctx context.Context) ([]model.ShopProduct, error) {
	var shopproducts []model.ShopProduct
	if err := r.db.WithContext(ctx).Find(&shopproducts).Error; err != nil {
		return nil, err
	}
	return shopproducts, nil
}

// Update updates an existing shopproduct
func (r *postgresShopRepository) ShopProductRepositoryUpdate(ctx context.Context, shopproduct *model.ShopProduct) error {
	if err := r.db.WithContext(ctx).Save(shopproduct).Error; err != nil {
		return err
	}
	return nil
}

// Delete removes a shopproduct from the database
func (r *postgresShopRepository) ShopProductRepositoryDelete(ctx context.Context, id int) error {
	if err := r.db.WithContext(ctx).Delete(&model.ShopProduct{}, id).Error; err != nil {
		log.Error(err)
		return err
	}
	return nil
}
