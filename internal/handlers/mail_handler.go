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

// MessageResponse defines a typical successful response for mail endpoints.
type MessageResponse struct {
	Message string `json:"message" example:"Message sent successfully"`
}

// ErrorResponse defines an error response.
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request"`
}

// SendMailRequest represents the payload for sending mail.
type SendMailRequest struct {
	ReceiverEmail string `json:"receiver_email" example:"receiver@example.com"`
	Subject       string `json:"subject" example:"Test Email"`
	Body          string `json:"body" example:"Hello, this is a test email."`
}

// @Summary Get mail service status
// @Description Returns the operational status of the mail service.
// @Produce json
// @Success 200 {object} MessageResponse "Mail service is operational"
// @Router /mail/status [get]
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

// @Summary Send mail
// @Description Send an email message.
// @Accept json
// @Produce json
// @Param message body SendMailRequest true "Mail message"
// @Success 200 {object} MessageResponse "Message sent successfully"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /mail/send [post]
// @Security BearerAuth
func (mh *MailHandler) SendMail(c *gin.Context) {
	claimsVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	claims := claimsVal.(*auth.Claims)

	var req SendMailRequest

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

// @Summary Get mail messages
// @Description Retrieve sent and received mail messages for the authenticated user.
// @Produce json
// @Success 200 {object} MailMessagesResponse "Mail messages"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /mail/messages [get]
// @Security BearerAuth
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

	c.JSON(http.StatusOK, gin.H{
		"sent":     sent,
		"received": received,
	})
}

// MailMessagesResponse defines the reply structure for GetMessages.
type MailMessagesResponse struct {
	Sent     []interface{} `json:"sent"`
	Received []interface{} `json:"received"`
}
