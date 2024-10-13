package service

import (
	"context"
	"simcomm-monolith/config"
	"simcomm-monolith/internal/model"
	"simcomm-monolith/internal/repository"
	"simcomm-monolith/util"

	log "github.com/labstack/gommon/log"
)

// ShopService defines the methods for the Shop service
type ShopService interface {
	Create(ctx context.Context, shop *model.Shop) error
	Get(ctx context.Context, id int) (*model.Shop, error)
	GetAll(ctx context.Context) ([]model.Shop, error)
	Update(ctx context.Context, shop *model.Shop) error
	Delete(ctx context.Context, id int) error
}

type shopService struct {
	repo      repository.ShopRepository
	redisRepo repository.RedisRepository
	cfg       *config.Config
}

func NewShopService(repo repository.ShopRepository, redisRepo repository.RedisRepository, cfg *config.Config) *shopService {
	return &shopService{
		repo:      repo,
		redisRepo: redisRepo,
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
