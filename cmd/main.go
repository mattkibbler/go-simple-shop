package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mattkibbler/go-simple-shop/internal/output"
	"github.com/mattkibbler/go-simple-shop/internal/shop"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func main() {
	listenAddr := ":8080"
	productCachePath := "product_cache.json"

	prodStore := shop.NewStore()
	templates := template.Must(template.New("").Funcs(template.FuncMap{
		"floatToCurrency": func(value float64) string {
			// Create a new printer for the English locale
			printer := message.NewPrinter(language.BritishEnglish)
			// Format the amount as USD
			return printer.Sprintf("£%.2f", value)
		},
	}).ParseGlob("internal/templates/*.html"))

	mux := mux.NewRouter()

	// Serve static assets
	staticDir := http.FileServer(http.Dir("public/assets"))
	mux.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", staticDir))
	// Web routes
	mux.HandleFunc("/", shop.HandleGetProducts(prodStore, templates)).Methods("GET")
	mux.HandleFunc("/product/{id}", shop.HandleGetProduct(prodStore, templates)).Methods("GET")
	// Handle 404s with custom page
	mux.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		output.WriteErrorPage(templates, w, errors.New("page not found"))
	})

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
