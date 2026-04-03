package models

type Wallet struct {
	ID string 
	UserID string
	Balance float64
}

type WalletResponse struct {
	UserID string
	Balance float64
}