package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Dav16Akin/payment-api/internal/models"
	"github.com/Dav16Akin/payment-api/internal/services"
)

type TransactionHandler interface {
	Transfer(w http.ResponseWriter, r *http.Request)
}

type transactionHandler struct {
	services services.TransactionService
}

func NewTransactionHandler(services services.TransactionService) TransactionHandler {
	return &transactionHandler{services: services}
}

func (s *transactionHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var req models.TransactionRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		defer r.Body.Close()

		transaction := models.Transaction{
			SenderID:   req.SenderID,
			RecieverID: req.RecieverID,
			Amount:     req.Amount,
			Status:     "Pending",
		}

		if err := s.services.Transfer(&transaction); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "transfer successful",
		})
	}else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
