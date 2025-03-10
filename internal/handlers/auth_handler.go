package handlers

import (
	"net/http"

	"mail-service/internal/auth"
	"mail-service/internal/repository"
	"mail-service/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// TokenResponse defines a successful login response.
type TokenResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.some.token.here"`
}

// LoginRequest represents the login payload.
type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"password123"`
}

type AuthHandler struct {
	repo   *repository.UserRepository
	logger *zap.Logger
}

func NewAuthHandler(repo *repository.UserRepository, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{repo: repo, logger: logger}
}

// @Summary Login user
// @Description Authenticate user and return a JWT token.
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "User credentials"
// @Success 200 {object} TokenResponse "Successful login"
// @Failure 400 {object} map[string]string "{"error": "Invalid request"}"
// @Failure 401 {object} map[string]string "{"error": "Invalid credentials"}"
// @Router /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

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
