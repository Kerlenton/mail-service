package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"mail-service/internal/handlers"
	"mail-service/internal/models"
	"mail-service/internal/repository"
	"mail-service/internal/router"
	"mail-service/internal/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupHandlerTest(t *testing.T) *gin.Engine {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory DB: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}
	repo := repository.NewUserRepository(db)
	logger := zap.NewNop()
	svc := services.NewUserService(repo, logger)
	handler := handlers.NewUserHandler(svc, logger)
	r := gin.Default()
	router.SetupRouter(r, handler)
	return r
}

func TestRegisterUser_Success(t *testing.T) {
	r := setupHandlerTest(t)
	requestBody, _ := json.Marshal(map[string]string{
		"email":    "test@example.com",
		"password": "secretpass",
	})
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestRegisterUser_InvalidJSON(t *testing.T) {
	r := setupHandlerTest(t)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte("{invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
