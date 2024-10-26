package shop

import (
	"fmt"
	"net/http"
)

func HandleGetProducts(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	fmt.Fprint(w, "products page...")
}

func HandleGetProduct(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	fmt.Fprint(w, "single products page...")
}
