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

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	Get(ctx context.Context, id int) (*model.User, error)
	GetAll(ctx context.Context) ([]model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int) error

	GetByIdentifier(ctx context.Context, identifier string) (*model.User, error)
}

type postgresUserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewPostgreUserRepository(db *gorm.DB) *postgresUserRepository {
	return &postgresUserRepository{db: db}
}

// Create inserts a new user into the database
func (r *postgresUserRepository) Create(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		if strings.Contains(err.Error(), util.SQLSTATE_23505) {
			return errors.New(util.ErrUserAlreadyExists)
		}
		log.Error(err)
		return errors.New(util.ErrInternalServerError)
	}

	return nil
}

// Get retrieves a user by ID
func (r *postgresUserRepository) Get(ctx context.Context, id int) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAll retrieves all users from the database
func (r *postgresUserRepository) GetAll(ctx context.Context) ([]model.User, error) {
	var users []model.User
	if err := r.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// Update updates an existing user
func (r *postgresUserRepository) Update(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(util.ErrUserNotFound)
		}
		log.Error(err)
		return err
	}

	return nil
}

// Delete removes a user from the database
func (r *postgresUserRepository) Delete(ctx context.Context, id int) error {
	if err := r.db.WithContext(ctx).Delete(&model.User{}, id).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(util.ErrUserNotFound)
		}

		log.Error(err)
		return err
	}
	return nil
}

func (r *postgresUserRepository) GetByIdentifier(ctx context.Context, identifier string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).
		Where("email = ? OR phone = ?", identifier, identifier).
		First(&user).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(util.ErrUserNotFound)
		}

		log.Error(err)
		return nil, err
	}
	return &user, nil
}
