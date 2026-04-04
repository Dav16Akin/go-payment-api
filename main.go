package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Dav16Akin/payment-api/internal/database"
	"github.com/Dav16Akin/payment-api/internal/handlers"
	"github.com/Dav16Akin/payment-api/internal/middleware"
	"github.com/Dav16Akin/payment-api/internal/repository"
	"github.com/Dav16Akin/payment-api/internal/services"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	db, err := database.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}

	if err := database.InitializeDB(db); err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	mux := http.NewServeMux()

	userRepo := repository.NewUserRepository(db)
	transactionRepo := repository.NewTransactionRepository()
	walletRepo := repository.NewWalletRepository(db)

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
