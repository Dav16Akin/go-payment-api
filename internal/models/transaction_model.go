package models

type Transaction struct {
	ID         string
	SenderID   string
	ReceiverID string
	Amount     float64
	Status     string
}

type TransactionRequest struct {
	SenderID   string  `json:"sender_id"`
	ReceiverID string  `json:"receiver_id"`
	Amount     float64 `json:"amount"`
}
