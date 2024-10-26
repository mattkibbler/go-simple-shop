package shop

import (
	"fmt"
	"net/http"
)

func HandleGetProducts(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprint(w, "products page...")
	}
}

func HandleGetProduct(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprint(w, "single products page...")
	}
}
