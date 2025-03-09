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

	// Initialize handlers.
	mailHandler := handlers.NewMailHandler()
	// For the admin handler, pass dummy nil logger and repo (for testing dummy responses).
	adminHandler := handlers.NewAdminHandler(nil, nil)
	router.SetupExpandedRoutes(r, mailHandler, adminHandler)

	// Test mail status endpoint.
	reqStatus, _ := http.NewRequest("GET", "/mail/status", nil)
	recStatus := httptest.NewRecorder()
	r.ServeHTTP(recStatus, reqStatus)
	assert.Equal(t, http.StatusOK, recStatus.Code)
	var statusResp map[string]string
	assert.NoError(t, json.Unmarshal(recStatus.Body.Bytes(), &statusResp))
	assert.Equal(t, "Mail service is operational", statusResp["status"])

	// Test mail history endpoint.
	reqHistory, _ := http.NewRequest("GET", "/mail/history", nil)
	recHistory := httptest.NewRecorder()
	r.ServeHTTP(recHistory, reqHistory)
	assert.Equal(t, http.StatusOK, recHistory.Code)
	var historyResp map[string][]map[string]interface{}
	assert.NoError(t, json.Unmarshal(recHistory.Body.Bytes(), &historyResp))
	assert.Len(t, historyResp["history"], 2)

	// Test admin endpoints.
	reqAdminUsers, _ := http.NewRequest("GET", "/admin/users", nil)
	recAdminUsers := httptest.NewRecorder()
	r.ServeHTTP(recAdminUsers, reqAdminUsers)
	assert.Equal(t, http.StatusOK, recAdminUsers.Code)
	var adminUsersResp map[string]string
	assert.NoError(t, json.Unmarshal(recAdminUsers.Body.Bytes(), &adminUsersResp))
	assert.Equal(t, "User list endpoint (not implemented)", adminUsersResp["message"])

	reqAdminMails, _ := http.NewRequest("POST", "/admin/mails/manage", nil)
	recAdminMails := httptest.NewRecorder()
	r.ServeHTTP(recAdminMails, reqAdminMails)
	assert.Equal(t, http.StatusOK, recAdminMails.Code)
	var adminMailsResp map[string]string
	assert.NoError(t, json.Unmarshal(recAdminMails.Body.Bytes(), &adminMailsResp))
	assert.Equal(t, "Mail management endpoint (not implemented)", adminMailsResp["message"])
}
