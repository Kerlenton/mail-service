package test

import (
	"context"
	"mail-service/internal/models"
	"mail-service/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.User{})
	return db, nil
}

func TestCreateUser(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	repo := repository.NewUserRepository(db)

	user := &models.User{Email: "test@example.com", PasswordHash: "hashedpassword"}
	err = repo.CreateUser(context.Background(), user)

	assert.NoError(t, err)
}
