package handlers

import (
	"net/http"

	"mail-service/internal/auth"
	"mail-service/internal/repository"
	"mail-service/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	repo   *repository.UserRepository
	logger *zap.Logger
}

func NewAuthHandler(repo *repository.UserRepository, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{repo: repo, logger: logger}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid login request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	user, err := h.repo.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil || user == nil {
		h.logger.Warn("User not found", zap.String("email", req.Email))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		h.logger.Warn("Invalid password", zap.String("email", req.Email))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	token, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		h.logger.Error("Failed to generate token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
