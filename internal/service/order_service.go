package service

import (
	"context"
	"simcomm-monolith/config"
	"simcomm-monolith/internal/model"
	"simcomm-monolith/internal/repository"
	"simcomm-monolith/util"

	log "github.com/labstack/gommon/log"
)

// OrderService defines the methods for the Order service
type OrderService interface {
	Create(ctx context.Context, order *model.Order) error
	Get(ctx context.Context, id int) (*model.Order, error)
	GetAll(ctx context.Context) ([]model.Order, error)
	Update(ctx context.Context, order *model.Order) error
	Delete(ctx context.Context, id int) error
}

type orderService struct {
	repo      repository.OrderRepository
	redisRepo repository.RedisRepository
	cfg       *config.Config
}

func NewOrderService(repo repository.OrderRepository, redisRepo repository.RedisRepository, cfg *config.Config) *orderService {
	return &orderService{
		repo:      repo,
		redisRepo: redisRepo,
		cfg:       cfg,
	}
}

func (s *orderService) Create(ctx context.Context, order *model.Order) error {
	timeNow := util.TimeNow()
	order.CreatedAt = timeNow
	order.UpdatedAt = timeNow
	err := s.repo.Create(ctx, order)
	if err != nil {
		log.Error(err)
	}
	return err
}

func (s *orderService) Get(ctx context.Context, id int) (*model.Order, error) {
	return s.repo.Get(ctx, id)
}

func (s *orderService) GetAll(ctx context.Context) ([]model.Order, error) {
	return s.repo.GetAll(ctx)
}

func (s *orderService) Update(ctx context.Context, order *model.Order) error {
	return s.repo.Update(ctx, order)
}

func (s *orderService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
