package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Dav16Akin/payment-api/internal/database"
	"github.com/Dav16Akin/payment-api/internal/handlers"
	"github.com/Dav16Akin/payment-api/internal/middleware"
	"github.com/Dav16Akin/payment-api/internal/repository"
	"github.com/Dav16Akin/payment-api/internal/services"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
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
	transactionRepo := repository.NewTransactionRepository(db)
	walletRepo := repository.NewWalletRepository(db)

	userService := services.NewUserService(userRepo, walletRepo)
	userHandler := handlers.NewUserHandler(userService)

	transactionService := services.NewTransactionService(walletRepo, transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	walletService := services.NewWalletService(walletRepo)
	walletHandler := handlers.NewWalletHandler(walletService)

	mux.HandleFunc("/sign-up", userHandler.SignUp)
	mux.HandleFunc("/sign-in", userHandler.SignIn)
	mux.HandleFunc("/transfer", transactionHandler.Transfer)
	mux.HandleFunc("/transactions", transactionHandler.GetAll)
	mux.HandleFunc("/transactions", transactionHandler.GetByUser)
	mux.HandleFunc("/wallet", walletHandler.GetWallet)

	handler := middleware.Logging(middleware.CORSMiddleware(mux))

	fmt.Println("Server running on PORT", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}
