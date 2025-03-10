package handlers

import (
	"context"
	"net/http"
	"time"

	"mail-service/internal/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RegisterRequest represents the user registration payload.
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserHandler struct {
	service *services.UserService
	logger  *zap.Logger
}

func NewUserHandler(service *services.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{service: service, logger: logger}
}

// @Summary Register new user
// @Description Create a new user account with email and password.
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "User credentials"
// @Success 201 {string} string "User created successfully"
// @Failure 400 {object} map[string]string "{"error": "Invalid request"}"
// @Router /register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var request RegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Warn("Invalid request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := h.service.RegisterUser(ctx, request.Email, request.Password); err != nil {
		h.logger.Error("Failed to register", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register"})
		return
	}

	c.Status(http.StatusCreated)
}
