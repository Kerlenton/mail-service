package test

import (
	"bytes"
	"encoding/json"
	"mail-service/internal/handlers"
	"mail-service/internal/models"
	"mail-service/internal/repository"
	"mail-service/internal/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRegisterHandler(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	db.AutoMigrate(&models.User{})

	repo := repository.NewUserRepository(db)
	service := services.NewUserService(repo)
	handler := handlers.NewUserHandler(service)

	requestBody, _ := json.Marshal(map[string]string{
		"email":    "test@example.com",
		"password": "securepassword",
	})

	r := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()
	handler.Register(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
}
