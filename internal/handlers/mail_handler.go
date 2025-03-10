package handlers

import (
	"net/http"

	"mail-service/internal/auth"
	"mail-service/internal/services"

	"github.com/gin-gonic/gin"
)

type MailHandler struct {
	messageService *services.MessageService
}

func NewMailHandler(messageService *services.MessageService) *MailHandler {
	return &MailHandler{messageService: messageService}
}

func (mh *MailHandler) GetMailStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "Mail service is operational"})
}

func (mh *MailHandler) GetMailHistory(c *gin.Context) {
	history := []gin.H{
		{"to": "user1@example.com", "subject": "Welcome", "sent_at": "2023-01-01T12:00:00Z"},
		{"to": "user2@example.com", "subject": "Notification", "sent_at": "2023-01-02T15:30:00Z"},
	}
	c.JSON(http.StatusOK, gin.H{"history": history})
}

func (mh *MailHandler) SendMail(c *gin.Context) {
	claimsVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	claims := claimsVal.(*auth.Claims)
	var req struct {
		ReceiverEmail string `json:"receiver_email"`
		Subject       string `json:"subject"`
		Body          string `json:"body"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if err := mh.messageService.SendMessage(claims.UserID, req.ReceiverEmail, req.Subject, req.Body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Message sent successfully"})
}

func (mh *MailHandler) GetMessages(c *gin.Context) {
	claimsVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	claims := claimsVal.(*auth.Claims)
	sent, received, err := mh.messageService.GetMessages(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"sent": sent, "received": received})
}
