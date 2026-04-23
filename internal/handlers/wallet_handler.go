package handlers

import (
	"net/http"

	"github.com/Dav16Akin/payment-api/internal/middleware"
	"github.com/Dav16Akin/payment-api/internal/models"
	"github.com/Dav16Akin/payment-api/internal/services"
	"github.com/Dav16Akin/payment-api/internal/utils"
)

type WalletHandler interface {
	GetWallet(w http.ResponseWriter, r *http.Request)
}

type walletHandler struct {
	services services.WalletService
}

func NewWalletHandler(services services.WalletService) WalletHandler {
	return &walletHandler{services: services}
}

func (h *walletHandler) GetWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.JSONResponse(w, http.StatusMethodNotAllowed, nil, "method not allowed")
		return
	}

	if r.Method == "GET" {
		userID, ok := r.Context().Value(middleware.UserIDKey).(string)
		if !ok {
			utils.JSONResponse(w, http.StatusUnauthorized, nil, "unauthorized")
			return
		}

		walletData, err := h.services.GetWallet(userID)
		if err != nil {
			utils.JSONResponse(w, http.StatusNotFound, nil, "wallet not found")
			return
		}

		wallet := models.WalletResponse{
			UserID:  walletData.UserID,
			Balance: walletData.Balance,
		}

		utils.JSONResponse(w, http.StatusOK, wallet, "")
	}
}
