package router

import (
	"mail-service/internal/handlers"
	"mail-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(r *gin.Engine, authHandler *handlers.AuthHandler) {
	// Public auth endpoint
	r.POST("/login", authHandler.Login)

	// Example group of protected routes
	protected := r.Group("/protected")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/hello", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Hello, authenticated user!"})
		})
	}
}
