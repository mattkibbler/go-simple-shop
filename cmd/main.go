package main

import (
	"net/http"

	"github.com/mattkibbler/go-simple-shop/internal/shop"
)

func main() {
	listenAddr := ":8080"

	mux := http.NewServeMux()

	mux.HandleFunc("/", shop.HandleGetProducts)
	mux.HandleFunc("GET /product/{id}", shop.HandleGetProduct)

	http.ListenAndServe(listenAddr, mux)
}
