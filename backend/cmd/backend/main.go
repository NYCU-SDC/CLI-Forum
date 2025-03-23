package main

import (
	"backend/internal/auth"
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /login", auth.LoginHandler)

	fmt.Println("Server up")

	http.ListenAndServe("localhost:8090", mux)
}
