package main

import (
	"backend/internal/auth"
	"fmt"
	"net/http"
)

func fooHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("foo handler")
	fmt.Fprintln(w, "foo")
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /login", auth.LoginHandler)
	mux.HandleFunc("POST /register", auth.RegisterHandler)
	mux.HandleFunc("/foo", auth.JWTMiddleware(fooHandler))

	fmt.Println("Server up")

	http.ListenAndServe("localhost:8090", mux)
}
