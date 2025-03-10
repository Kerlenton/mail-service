package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"mail-service/internal/handlers"
	"mail-service/internal/router"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestExpandedRouterEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mailHandler := handlers.NewMailHandler(nil)
	adminHandler := handlers.NewAdminHandler(nil, nil)
	router.SetupExpandedRoutes(r, mailHandler, adminHandler)
	reqStatus, _ := http.NewRequest("GET", "/mail/status", nil)
	recStatus := httptest.NewRecorder()
	r.ServeHTTP(recStatus, reqStatus)
	assert.Equal(t, http.StatusOK, recStatus.Code)
	var statusResp map[string]string
	assert.NoError(t, json.Unmarshal(recStatus.Body.Bytes(), &statusResp))
	assert.Equal(t, "Mail service is operational", statusResp["status"])
	reqMessages, _ := http.NewRequest("GET", "/mail/messages", nil)
	recMessages := httptest.NewRecorder()
	r.ServeHTTP(recMessages, reqMessages)
	assert.Equal(t, http.StatusOK, recMessages.Code)
	var messagesResp map[string][]interface{}
	assert.NoError(t, json.Unmarshal(recMessages.Body.Bytes(), &messagesResp))
	assert.Empty(t, messagesResp["sent"])
	assert.Empty(t, messagesResp["received"])
}
