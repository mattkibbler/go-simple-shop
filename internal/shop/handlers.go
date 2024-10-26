package shop

import (
	"html/template"
	"net/http"

	"github.com/mattkibbler/go-simple-shop/internal/output"
)

func HandleGetProducts(store *Store, templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		pageData := struct {
			Title string
		}{
			Title: "Products",
		}
		output.RenderPage(templates, w, "products.html", pageData)
	}
}

func HandleGetProduct(store *Store, templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		pageData := struct {
			Title string
		}{
			Title: "Products",
		}
		output.RenderPage(templates, w, "product.html", pageData)
	}
}
