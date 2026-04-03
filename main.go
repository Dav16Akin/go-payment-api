package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Dav16Akin/payment-api/internal/handlers"
	"github.com/Dav16Akin/payment-api/internal/middleware"
	"github.com/Dav16Akin/payment-api/internal/models"
	"github.com/Dav16Akin/payment-api/internal/repository"
	"github.com/Dav16Akin/payment-api/internal/services"
)

func main() {
	mux := http.NewServeMux()
	userRepo := repository.NewUserRepository()
	transactionRepo := repository.NewTransactionRepository()
	walletRepo := repository.NewWalletRepository()


	userRepo.CreateUser(&models.User{
		ID:    "user1",
		Name:  "David",
		Email: "david@test.com",
	})

	userRepo.CreateUser(&models.User{
		ID:    "user2",
		Name:  "John",
		Email: "john@test.com",
	})


	walletRepo.CreateWallet(&models.Wallet{
		ID:      "wallet1",
		UserID:  "user1",
		Balance: 1000,
	})

	walletRepo.CreateWallet(&models.Wallet{
		ID:      "wallet2",
		UserID:  "user2",
		Balance: 500,
	})

	userService := services.NewUserService(userRepo, walletRepo)
	userHandler := handlers.NewUserHandler(userService)

	transactionService := services.NewTransactionService(walletRepo, transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	walletService := services.NewWalletService(walletRepo)
	walletHandler := handlers.NewWalletHandler(walletService)

	mux.HandleFunc("/user", userHandler.CreateUser)
	mux.HandleFunc("/transfer", transactionHandler.Transfer)
	mux.HandleFunc("/transactions", transactionHandler.GetAll)
	mux.HandleFunc("/wallet/{user_id}", walletHandler.GetWallet)


	loggedMux := middleware.Logging(mux)

	fmt.Println("Server running on PORT 8000")
	if err := http.ListenAndServe(":8000", loggedMux); err != nil {
		log.Fatal(err)
	}
}
