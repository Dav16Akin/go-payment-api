package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Dav16Akin/payment-api/internal/handlers"
	"github.com/Dav16Akin/payment-api/internal/models"
	"github.com/Dav16Akin/payment-api/internal/repository"
	"github.com/Dav16Akin/payment-api/internal/services"
)

func main() {
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

	userRepo.ListAll()
	walletRepo.ListAllWallets()

	userService := services.NewUserService(userRepo, walletRepo)
	userHandler := handlers.NewUserHandler(userService)

	transactionService := services.NewTransactionService(walletRepo, transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	walletService := services.NewWalletService(walletRepo)
	walletHandler := handlers.NewWalletHandler(walletService)

	http.HandleFunc("/user", userHandler.CreateUser)
	http.HandleFunc("/transfer", transactionHandler.Transfer)
	http.HandleFunc("/transactions", transactionHandler.GetAll)
	http.HandleFunc("/wallet/{user_id}", walletHandler.GetWallet)


	fmt.Println("Server running on PORT 8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
