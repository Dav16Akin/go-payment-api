package models

type Transaction struct {
	ID         string
	SenderID   string
	RecieverID string
	Amount     float64
	Status     string
}

type TransactionRequest struct {
	SenderID   string  `json:"sender_id"`
	RecieverID string  `json:"receiver_id"`
	Amount     float64 `json:"amount"`
}
