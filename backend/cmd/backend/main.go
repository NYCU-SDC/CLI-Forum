package main

import (
	"backend/internal/auth"
	"backend/internal/database"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file : ", err)
	}

	database.MigrateUP()
	defer database.MigrateDown()

	// initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic("Failed to initialize Zap logger: " + err.Error())
	}

	defer func() {
		err := logger.Sync()
		if err != nil {
			fmt.Println("Failed to sync logger: ", err)
		}
	}() // flushes buffer, if any

	sugar := logger.Sugar()

	// initialize database
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic("Unable to create connection pool: " + err.Error())
	}
	defer dbpool.Close()

	// initialize auth service
	authService := auth.NewService(dbpool)
	authHandler := auth.NewHandler(authService)

	// initialize mux
	mux := http.NewServeMux()

	// set up routes
	mux.HandleFunc("POST /login", authHandler.LoginHandler)
	mux.HandleFunc("POST /register", authHandler.RegisterHandler)

	// start server on port 8090
	sugar.Info("Server starting on localhost:8090")
	err = http.ListenAndServe("localhost:8090", mux)
	if err != nil {
		sugar.Fatal("Fail to start server with error : ", err)
	}
}
