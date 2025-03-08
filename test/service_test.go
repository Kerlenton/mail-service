package test

import (
	"mail-service/internal/models"
	"mail-service/internal/repository"
	"mail-service/internal/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUserService(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	db.AutoMigrate(&models.User{})
	repo := repository.NewUserRepository(db)
	service := services.NewUserService(repo)

	email := "test@example.com"
	password := "securepassword"
	assert.NoError(t, service.RegisterUser(email, password))

	fetchedUser, err := service.GetUserByEmail(email)
	assert.NoError(t, err)
	assert.Equal(t, email, fetchedUser.Email)
}
