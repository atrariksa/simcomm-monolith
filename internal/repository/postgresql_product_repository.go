package repository

import (
	"context"
	"errors"
	"simcomm-monolith/internal/model"
	"simcomm-monolith/util"
	"strings"

	log "github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) error
	Get(ctx context.Context, id int) (*model.Product, error)
	GetAll(ctx context.Context) ([]model.Product, error)
	Update(ctx context.Context, product *model.Product) error
	Delete(ctx context.Context, id int) error
}

type postgresProductRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new instance of ProductRepository
func NewPostgreProductRepository(db *gorm.DB) *postgresProductRepository {
	return &postgresProductRepository{db: db}
}

// Create inserts a new product into the database
func (r *postgresProductRepository) Create(ctx context.Context, product *model.Product) error {
	if err := r.db.WithContext(ctx).Create(product).Error; err != nil {
		log.Error(err)
		if strings.Contains(err.Error(), util.SQLSTATE_23505) {
			return errors.New(util.ErrUserAlreadyExists)
		}
		log.Error(err)
		return errors.New(util.ErrInternalServerError)
	}

	return nil
}

// Get retrieves a product by ID
func (r *postgresProductRepository) Get(ctx context.Context, id int) (*model.Product, error) {
	var product model.Product
	if err := r.db.WithContext(ctx).First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

// GetAll retrieves all products from the database
func (r *postgresProductRepository) GetAll(ctx context.Context) ([]model.Product, error) {
	var products []model.Product
	if err := r.db.WithContext(ctx).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// Update updates an existing product
func (r *postgresProductRepository) Update(ctx context.Context, product *model.Product) error {
	if err := r.db.WithContext(ctx).Save(product).Error; err != nil {
		return err
	}
	return nil
}

// Delete removes a product from the database
func (r *postgresProductRepository) Delete(ctx context.Context, id int) error {
	if err := r.db.WithContext(ctx).Delete(&model.Product{}, id).Error; err != nil {
		return err
	}
	return nil
}
