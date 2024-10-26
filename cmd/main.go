package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/mattkibbler/go-simple-shop/internal/shop"
)

var prodCache sync.Map

func main() {
	listenAddr := ":8080"

	mux := http.NewServeMux()

	mux.HandleFunc("/", shop.HandleGetProducts)
	mux.HandleFunc("GET /product/{id}", shop.HandleGetProduct)

	go func() {
		for {
			log.Println("Fetching products...")
			err := shop.FetchAndCacheProducts(&prodCache)
			if err != nil {
				log.Printf("error fetching products: %v\n", err)
			}
			// Schedule product refresh after a delay
			time.Sleep(10 * time.Minute)
		}
	}()

	log.Printf("Server listening at %v", listenAddr)
	http.ListenAndServe(listenAddr, mux)
}
