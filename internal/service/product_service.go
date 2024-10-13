package service

import (
	"context"
	"simcomm-monolith/config"
	"simcomm-monolith/internal/model"
	"simcomm-monolith/internal/repository"
	"simcomm-monolith/util"

	log "github.com/labstack/gommon/log"
)

// ProductService defines the methods for the Product service
type ProductService interface {
	Create(ctx context.Context, product *model.Product) error
	Get(ctx context.Context, id int) (*model.Product, error)
	GetAll(ctx context.Context) ([]model.Product, error)
	Update(ctx context.Context, product *model.Product) error
	Delete(ctx context.Context, id int) error
}

type productService struct {
	repo      repository.ProductRepository
	redisRepo repository.RedisRepository
	cfg       *config.Config
}

func NewProductService(repo repository.ProductRepository, redisRepo repository.RedisRepository, cfg *config.Config) *productService {
	return &productService{
		repo:      repo,
		redisRepo: redisRepo,
		cfg:       cfg,
	}
}

func (s *productService) Create(ctx context.Context, product *model.Product) error {
	timeNow := util.TimeNow()
	product.CreatedAt = timeNow
	product.UpdatedAt = timeNow
	err := s.repo.Create(ctx, product)
	if err != nil {
		log.Error(err)
	}
	return err
}

func (s *productService) Get(ctx context.Context, id int) (*model.Product, error) {
	return s.repo.Get(ctx, id)
}

func (s *productService) GetAll(ctx context.Context) ([]model.Product, error) {
	return s.repo.GetAll(ctx)
}

func (s *productService) Update(ctx context.Context, product *model.Product) error {
	timeNow := util.TimeNow()
	product.UpdatedAt = timeNow
	return s.repo.Update(ctx, product)
}

func (s *productService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
