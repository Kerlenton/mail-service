package router

import (
	"mail-service/internal/handlers"
	"mail-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupExpandedRoutes(r *gin.Engine, mailHandler *handlers.MailHandler, adminHandler *handlers.AdminHandler) {
	// Public endpoint
	r.GET("/mail/status", mailHandler.GetMailStatus)

	// Protected mail routes
	protectedMail := r.Group("/mail").Use(middleware.AuthMiddleware())
	{
		protectedMail.POST("/send", mailHandler.SendMail)
		protectedMail.GET("/messages", mailHandler.GetMessages)
	}
}
