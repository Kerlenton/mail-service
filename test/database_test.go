package test

import (
	"testing"

	"mail-service/internal/models"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupInMemoryDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory DB: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("migration failed: %v", err)
	}
	return db
}

func TestDatabase_MigrationAndCRUD(t *testing.T) {
	db := setupInMemoryDB(t)
	user := models.User{Email: "mock@example.com", PasswordHash: "hashed"}
	err := db.Create(&user).Error
	assert.NoError(t, err, "should create user without error")
	var found models.User
	err = db.First(&found, "email = ?", "mock@example.com").Error
	assert.NoError(t, err, "should find the user")
	assert.Equal(t, "mock@example.com", found.Email)
}
