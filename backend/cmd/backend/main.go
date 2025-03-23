package main

import (
	"backend/internal/auth"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func main() {
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

	// initialize mux
	mux := http.NewServeMux()

	// set up routes
	mux.HandleFunc("POST /login", auth.LoginHandler)
	mux.HandleFunc("POST /register", auth.RegisterHandler)

	// start server on port 8090
	sugar.Info("Server starting on localhost:8090")
	err = http.ListenAndServe("localhost:8090", mux)
	if err != nil {
		sugar.Fatal("Fail to start server with error : ", err)
	}
}
