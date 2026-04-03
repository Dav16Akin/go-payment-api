package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Dav16Akin/payment-api/internal/models"
	"github.com/Dav16Akin/payment-api/internal/services"
	"github.com/Dav16Akin/payment-api/internal/utils"
)

type TransactionHandler interface {
	Transfer(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
}

type transactionHandler struct {
	services services.TransactionService
}

func NewTransactionHandler(services services.TransactionService) TransactionHandler {
	return &transactionHandler{services: services}
}

func (h *transactionHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.JSONResponse(w, http.StatusMethodNotAllowed, nil, "method not allowed")
		return
	}

	if r.Method == "POST" {
		var req models.TransactionRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			utils.JSONResponse(w, http.StatusBadRequest, nil, "invalid request body")
			return
		}

		defer r.Body.Close()

		transaction := models.Transaction{
			SenderID:   req.SenderID,
			ReceiverID: req.ReceiverID,
			Amount:     req.Amount,
			Status:     "pending",
		}

		if err := h.services.Transfer(&transaction); err != nil {
			utils.JSONResponse(w, http.StatusBadRequest, nil, err.Error())
			return
		}

		utils.JSONResponse(w, http.StatusCreated, map[string]string{
			"message": "transfer successful",
		}, "")
	}
}

func (h *transactionHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.JSONResponse(w, http.StatusMethodNotAllowed, nil, "method not allowed")
		return
	}

	if r.Method == "GET" {
		transactions, err := h.services.GetAll()
		if err != nil {
			utils.JSONResponse(w, http.StatusInternalServerError, nil, "cannot get transactions")
			return
		}

		utils.JSONResponse(w, http.StatusOK, transactions, "")
	}
}
