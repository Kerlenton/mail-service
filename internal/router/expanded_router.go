package router

import (
	"mail-service/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupExpandedRoutes(r *gin.Engine, mailHandler *handlers.MailHandler, adminHandler *handlers.AdminHandler) {
	// Mail endpoints
	r.GET("/mail/status", mailHandler.GetMailStatus)
	r.GET("/mail/history", mailHandler.GetMailHistory)

	// Admin endpoints
	admin := r.Group("/admin")
	{
		admin.GET("/users", adminHandler.ListUsers)
		admin.POST("/mails/manage", adminHandler.ManageMail)
	}
}
