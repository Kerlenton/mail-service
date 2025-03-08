package test

import (
	"mail-service/internal/models"
	"mail-service/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUserRepository(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	db.AutoMigrate(&models.User{})

	repo := repository.NewUserRepository(db)
	user := &models.User{Email: "test@example.com", PasswordHash: "hashedpassword"}
	assert.NoError(t, repo.CreateUser(user))

	fetchedUser, err := repo.GetUserByEmail("test@example.com")
	assert.NoError(t, err)
	assert.Equal(t, user.Email, fetchedUser.Email)
}
