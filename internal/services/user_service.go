package services

import (
	"context"
	"errors"
	"mail-service/internal/models"
	"mail-service/internal/repository"
	"mail-service/internal/utils"

	"go.uber.org/zap"
)

type UserService struct {
	repo   *repository.UserRepository
	logger *zap.Logger
}

func NewUserService(repo *repository.UserRepository, logger *zap.Logger) *UserService {
	return &UserService{repo: repo, logger: logger}
}

func (s *UserService) RegisterUser(ctx context.Context, email, password string) error {
	// Debug: log the repository pointer value
	s.logger.Debug("RegisterUser called", zap.Any("repo", s.repo))
	if len(password) < 6 {
		return errors.New("password too short")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		s.logger.Error("Failed to hash password", zap.Error(err))
		return err
	}

	user := &models.User{Email: email, PasswordHash: hashedPassword}
	if err := s.repo.CreateUser(ctx, user); err != nil {
		s.logger.Error("Failed to create user", zap.Error(err))
		return err
	}

	s.logger.Info("User registered successfully", zap.String("email", email))
	return nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		s.logger.Error("Failed to get user by email", zap.String("email", email), zap.Error(err))
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}
