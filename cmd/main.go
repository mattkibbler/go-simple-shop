package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mattkibbler/go-simple-shop/internal/output"
	"github.com/mattkibbler/go-simple-shop/internal/shop"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func main() {
	port := flag.String("port", "8080", "Port to run the application on")
	// Parse flags
	flag.Parse()
	listenAddr := fmt.Sprintf(":%v", *port)
	productCachePath := "product_cache.json"

	prodStore := shop.NewStore()
	templates := template.Must(template.New("").Funcs(template.FuncMap{
		"floatToCurrency": func(value float64) string {
			// Create a new printer for the English locale
			printer := message.NewPrinter(language.BritishEnglish)
			// Format the amount as USD
			return printer.Sprintf("£%.2f", value)
		},
		"ratingStars": func(value float64) []bool {
			// Round down rating to display stars...
			// A "4.51" rated product should probably be given 4 stars not 5
			rounded := int(math.Floor(value))
			var result []bool
			for i := 1; i <= 5; i++ {
				if i <= rounded {
					result = append(result, true)
				} else {
					result = append(result, false)
				}
			}
			return result
		},
		"formatDate": func(t time.Time) string {
			return t.Format("January 2, 2006")
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
		sleepDuration := 1 * time.Minute
		log.Println("Attempting to unserialize product cache...")
		pCount, err := shop.UnserializeProductCache(productCachePath, &prodStore.Cache)
		log.Printf("%d products unserialized\n", pCount)
		if err != nil {
			log.Printf("error unserializing product cache: %v\n", err)
		}
		// If we unserialized some products, we don't need to fetch them from the API straight away
		// This is mainly so we're not hitting the API every time we stop/start the server during development
		if pCount > 0 {
			time.Sleep(sleepDuration)
		}
		for {
			log.Println("Fetching products...")
			fetchedCount, err := shop.FetchAndCacheProducts(&prodStore.Cache)
			if err != nil {
				log.Printf("error fetching products: %v\n", err)
			}
			log.Printf("Fetched %v products", fetchedCount)
			log.Println("Serializing products...")
			err = shop.SerializeProductCache(productCachePath, &prodStore.Cache)
			if err != nil {
				log.Printf("error serializing product cache: %v\n", err)
			}
			// Schedule product refresh after a delay
			time.Sleep(sleepDuration)
		}
	}()

	log.Printf("Server listening at %v", listenAddr)
	http.ListenAndServe(listenAddr, mux)
}
