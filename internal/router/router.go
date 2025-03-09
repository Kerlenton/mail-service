package router

import (
	"mail-service/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine, userHandler *handlers.UserHandler) {
	r.POST("/register", userHandler.Register)
}
