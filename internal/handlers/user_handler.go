package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Dav16Akin/payment-api/internal/models"
	"github.com/Dav16Akin/payment-api/internal/services"
)

type UserHandler interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) UserHandler {
	return &userHandler{service: service}
}

func (h *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var req models.CreateUserRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			if strings.Contains(err.Error(), "required") {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		defer r.Body.Close()

		user := models.User{
			Name:  strings.TrimSpace(req.Name),
			Email: strings.TrimSpace(req.Email),
		}

		createdUser, err := h.service.CreateUser(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp := models.UserResponse{
			ID:    createdUser.ID,
			Name:  createdUser.Name,
			Email: createdUser.Email,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
