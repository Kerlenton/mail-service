package handlers

import (
	"net/http"

	"mail-service/internal/repository"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AdminHandler struct {
	repo   *repository.UserRepository
	logger *zap.Logger
}

func NewAdminHandler(repo *repository.UserRepository, logger *zap.Logger) *AdminHandler {
	return &AdminHandler{repo: repo, logger: logger}
}

func (ah *AdminHandler) ListUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "User list endpoint (not implemented)"})
}

func (ah *AdminHandler) ManageMail(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Mail management endpoint (not implemented)"})
}
