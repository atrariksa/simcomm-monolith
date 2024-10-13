package service

import (
	"context"
	"errors"
	"simcomm-monolith/config"
	"simcomm-monolith/internal/model"
	"simcomm-monolith/internal/repository"
	"simcomm-monolith/util"

	log "github.com/labstack/gommon/log"
)

// UserService defines the methods for the User service
type UserService interface {
	Create(ctx context.Context, user *model.User) error
	Get(ctx context.Context, id int) (*model.User, error)
	GetAll(ctx context.Context) ([]model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int) error

	GetUserByIdentifier(ctx context.Context, identifier string) (*model.User, error)
	SignUp(ctx context.Context, req model.SignUpRequest) error
	Login(ctx context.Context, req model.LoginRequest) (model.LoginData, error)
}

type userService struct {
	repo      repository.UserRepository
	redisRepo repository.RedisRepository
	cfg       *config.Config
}

func NewUserService(repo repository.UserRepository, redisRepo repository.RedisRepository, cfg *config.Config) *userService {
	return &userService{
		repo:      repo,
		redisRepo: redisRepo,
		cfg:       cfg,
	}
}

func (s *userService) Create(ctx context.Context, user *model.User) error {
	return s.repo.Create(ctx, user)
}

func (s *userService) Get(ctx context.Context, id int) (*model.User, error) {
	return s.repo.Get(ctx, id)
}

func (s *userService) GetAll(ctx context.Context) ([]model.User, error) {
	return s.repo.GetAll(ctx)
}

func (s *userService) Update(ctx context.Context, user *model.User) error {
	return s.repo.Update(ctx, user)
}

func (s *userService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *userService) GetUserByIdentifier(ctx context.Context, identifier string) (*model.User, error) {
	return s.repo.GetByIdentifier(ctx, identifier)
}

func (s *userService) SignUp(ctx context.Context, req model.SignUpRequest) error {

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		log.Error(err)
		return errors.New(util.ErrInternalServerError)
	}

	user := &model.User{
		Name:      req.Name,
		Phone:     req.Phone,
		Email:     req.Email,
		Passsword: hashedPassword,
		UserDetail: model.UserDetail{
			Roles: []string{req.Role},
		},
	}
	err = s.repo.Create(ctx, user)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (s *userService) Login(ctx context.Context, req model.LoginRequest) (model.LoginData, error) {
	loginData := model.LoginData{}
	user, err := s.GetUserByIdentifier(ctx, req.Identifier)
	if err != nil {
		log.Error(err)
		return loginData, err
	}

	cfg := s.cfg.AuthTokenConfig
	token, err := util.GenerateToken(user.ID, cfg.Duration, cfg.SecretKey)
	if err != nil {
		return loginData, err
	}
	loginData.Token = token

	return loginData, nil
}
