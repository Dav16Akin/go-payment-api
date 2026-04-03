package handlers

import (
	"net/http"

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
		id := r.PathValue("user_id")

		if id == "" {
			utils.JSONResponse(w, http.StatusBadRequest, nil, "user_id is required")
			return
		}

		walletData, err := h.services.GetWallet(id)
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
