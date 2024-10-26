package main

import (
	"log"
	"net/http"
	"time"

	"github.com/mattkibbler/go-simple-shop/internal/shop"
)

func main() {
	listenAddr := ":8080"
	productCachePath := "product_cache.json"

	prodStore := shop.NewStore()

	mux := http.NewServeMux()

	mux.HandleFunc("/", shop.HandleGetProducts)
	mux.HandleFunc("GET /product/{id}", shop.HandleGetProduct)

	go func() {
		log.Println("Attempting to unserialize product cache...")
		pCount, err := shop.UnserializeProductCache(productCachePath, &prodStore.Cache)
		log.Printf("%d products unserialized\n", pCount)
		if err != nil {
			log.Printf("error unserializing product cache: %v\n", err)
		}
		// If we unserialized some products, we don't need to fetch them from the API straight away
		// This is mainly so we're not hitting the API every time we stop/start the server during development
		if pCount > 0 {
			time.Sleep(10 * time.Minute)
		}
		for {
			log.Println("Fetching products...")
			err := shop.FetchAndCacheProducts(&prodStore.Cache)
			if err != nil {
				log.Printf("error fetching products: %v\n", err)
			}
			err = shop.SerializeProductCache(productCachePath, &prodStore.Cache)
			if err != nil {
				log.Printf("error serializing product cache: %v\n", err)
			}

			// Schedule product refresh after a delay
			time.Sleep(10 * time.Minute)
		}
	}()

	log.Printf("Server listening at %v", listenAddr)
	http.ListenAndServe(listenAddr, mux)
}
