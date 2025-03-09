package repository

import (
	"context"
	"errors"
	"mail-service/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	if db == nil {
		panic("NewUserRepository received a nil *gorm.DB")
	}
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	// Added nil check at the very start.
	if r.db == nil {
		return errors.New("gorm.DB pointer is nil in CreateUser")
	}
	// Check that the underlying SQL connection is active.
	sqlDB, err := r.db.DB()
	if err != nil || sqlDB == nil {
		return errors.New("underlying SQL connection is nil")
	}
	if err := sqlDB.Ping(); err != nil {
		return errors.New("underlying SQL connection is not active")
	}
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
