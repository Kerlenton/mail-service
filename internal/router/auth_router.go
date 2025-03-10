package router

import (
	"mail-service/internal/handlers"
	"mail-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(r *gin.Engine, authHandler *handlers.AuthHandler) {
	r.POST("/login", authHandler.Login)

	protected := r.Group("/protected")
	protected.Use(middleware.AuthMiddleware())
	{
		// @Summary Protected hello endpoint
		// @Description Returns a hello message for authenticated users.
		// @Produce json
		// @Success 200 {object} map[string]string "Hello, authenticated user!"
		// @Failure 401 {object} map[string]string "Unauthorized"
		// @Router /protected/hello [get]
		// @Security BearerAuth
		protected.GET("/hello", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Hello, authenticated user!"})
		})
	}
}
