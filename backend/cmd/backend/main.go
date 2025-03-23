package main

import (
	"backend/internal/auth"
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /login", auth.LoginHandler)
	mux.HandleFunc("POST /register", auth.RegisterHandler)

	fmt.Println("Server up")

	http.ListenAndServe("localhost:8090", mux)
}
