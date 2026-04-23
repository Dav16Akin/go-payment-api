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
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Println("Starting server...")
	log.Println("PORT:", port)
	log.Println("DATABASE_PUBLIC_URL set:", os.Getenv("DATABASE_PUBLIC_URL") != "")

	db, err := database.ConnectToDB()
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer db.Close()
	log.Println("DB connected successfully")

	if err := database.InitializeDB(db); err != nil {
		log.Fatal("Failed to initialize DB:", err)
	}
	log.Println("DB initialized successfully")

	err = database.RunMigrations(db)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	userRepo := repository.NewUserRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	tokenRepo := repository.NewTokenRepository(db)

	userService := services.NewUserService(userRepo, walletRepo, tokenRepo)
	userHandler := handlers.NewUserHandler(userService)

	transactionService := services.NewTransactionService(walletRepo, transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	walletService := services.NewWalletService(walletRepo)
	walletHandler := handlers.NewWalletHandler(walletService)

	protected := middleware.AuthMiddleware

	protectedHandler := func(h http.HandlerFunc) http.Handler {
		return protected(h)
	}


	//AUTHENTICATION
	mux.HandleFunc("/sign-up", userHandler.SignUp)
	mux.HandleFunc("/sign-in", userHandler.SignIn)
	mux.HandleFunc("/refresh", userHandler.RefreshToken)

	//TRANSACTIONS AND WALLETS
	mux.Handle("/transfer", protectedHandler(transactionHandler.Transfer))
	mux.Handle("/transactions", protectedHandler(transactionHandler.GetAll))
	mux.Handle("/transactions/user", protectedHandler(transactionHandler.GetByUser))
	mux.Handle("/wallet", protectedHandler(walletHandler.GetWallet))

	//USERS
	mux.Handle("/users/profile", protectedHandler(userHandler.UpdateProfile))
	mux.Handle("/users/password", protectedHandler(userHandler.ChangePassword))

	handler := middleware.Logging(middleware.CORSMiddleware(mux))

	fmt.Println("Server running on PORT", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}
