package router

import (
	"mail-service/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userHandler *handlers.UserHandler) *gin.Engine {
	r := gin.Default()
	r.POST("/register", userHandler.Register)
	return r
}
