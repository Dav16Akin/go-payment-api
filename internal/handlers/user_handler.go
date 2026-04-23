package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Dav16Akin/payment-api/internal/middleware"
	"github.com/Dav16Akin/payment-api/internal/models"
	"github.com/Dav16Akin/payment-api/internal/services"
	"github.com/Dav16Akin/payment-api/internal/utils"
)

type UserHandler interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	SignIn(w http.ResponseWriter, r *http.Request)

	GetUserProfile(w http.ResponseWriter, r *http.Request)
	UpdateUserProfile(w http.ResponseWriter, r *http.Request)
	ChangeUserPassword(w http.ResponseWriter, r *http.Request)

	RefreshToken(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) UserHandler {
	return &userHandler{service: service}
}

func (h *userHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.JSONResponse(w, http.StatusMethodNotAllowed, nil, "method not allowed")
		return
	}

	if r.Method == "POST" {
		var req models.CreateUserRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			if strings.Contains(err.Error(), "required") {
				utils.JSONResponse(w, http.StatusBadRequest, nil, "invalid request body")
			} else {
				utils.JSONResponse(w, http.StatusInternalServerError, nil, err.Error())
			}
			return
		}

		defer r.Body.Close()

		user := models.User{
			Name:     strings.TrimSpace(req.Name),
			Email:    strings.TrimSpace(req.Email),
			Password: strings.TrimSpace(req.Password),
		}

		createdUser, err := h.service.SignUp(&user)
		if err != nil {
			utils.JSONResponse(w, http.StatusBadRequest, nil, err.Error())
			return
		}

		resp := models.UserResponse{
			ID:    createdUser.ID,
			Name:  createdUser.Name,
			Email: createdUser.Email,
		}

		utils.JSONResponse(w, http.StatusCreated, resp, "")
	}
}

func (h *userHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.JSONResponse(w, http.StatusMethodNotAllowed, nil, "method not allowed")
		return
	}

	var req models.SignInRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, nil, "invalid requeest body")
		return
	}

	accessToken, refreshToken, err := h.service.SignIn(&req)
	if err != nil {
		utils.JSONResponse(w, http.StatusUnauthorized, nil, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		MaxAge:   7 * 24 * 60 * 60,
	})

	utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"accessToken": accessToken,
	}, "")
}

func (h *userHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.JSONResponse(w, http.StatusMethodNotAllowed, nil, "method not allowed")
		return
	}

	cookie, err := r.Cookie("refreshToken")
	if err != nil {
		if err == http.ErrNoCookie {
			utils.JSONResponse(w, http.StatusBadRequest, nil, "refresh token not found")
			return
		}
		utils.JSONResponse(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	accessToken, refreshToken, err := h.service.RefreshToken(cookie.Value)
	if err != nil {
		utils.JSONResponse(w, http.StatusUnauthorized, nil, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		MaxAge:   7 * 24 * 60 * 60,
	})

	utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"accessToken": accessToken,
	}, "")
}

func (h *userHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.JSONResponse(w, http.StatusMethodNotAllowed, nil, "method not allowed")
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		utils.JSONResponse(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}

	user, err := h.service.GetUserProfile(userID)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	avatarURL := ""
	if user.AvatarURL != nil {
		avatarURL = *user.AvatarURL
	}

	utils.JSONResponse(w, http.StatusOK, map[string]interface{}{"user": map[string]string{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"avatar_url": avatarURL,
	}}, "")
}

func (h *userHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPatch {
		utils.JSONResponse(w, http.StatusMethodNotAllowed, nil, "method not allowed")
		return
	}

	var req models.UpdateProfileRequest

	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		utils.JSONResponse(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, nil, "invalid request body")
		return
	}

	user, err := h.service.UpdateUserProfile(userID, &req)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	response := models.UpdateProfileResponse{
		ID:        user.ID,
		Name:      user.Name,
		AvatarURL: user.AvatarURL,
	}

	utils.JSONResponse(w, http.StatusCreated, response, "")

}

func (h *userHandler) ChangeUserPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		utils.JSONResponse(w, http.StatusMethodNotAllowed, nil, "method not allowed")
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		utils.JSONResponse(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}

	var req models.ChangePasswordRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, nil, "invalid request body")
		return
	}

	msg, err := h.service.ChangeUserPassword(userID, &req)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusOK, msg, "")

}
