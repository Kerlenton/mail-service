package handlers

import (
	"encoding/json"
	"net/http"

	"mail-service/internal/services"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.service.RegisterUser(request.Email, request.Password); err != nil {
		http.Error(w, "Failed to register", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
