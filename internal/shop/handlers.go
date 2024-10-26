package shop

import (
	"bytes"
	"html/template"
	"net/http"
	"strconv"

	"github.com/mattkibbler/go-simple-shop/internal/output"
)

func HandleGetProducts(store *Store, templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var buffer bytes.Buffer

		queryParams := r.URL.Query()
		pageStr := queryParams.Get("page") // Returns the first value for the key "name"
		page, _ := strconv.Atoi(pageStr)
		if page == 0 {
			page = 1
		}

		products, err := store.QueryProducts(func(p Product) bool {
			return true
		}, func(i Product, j Product) bool {
			return i.Title < j.Title
		})
		if err != nil {
			return
		}

		pageData := output.PageData{
			Title: "Products",
			Data: ProductsPageData{
				PaginatedData: *output.NewPaginatedPage(products, 12, page),
			},
		}

		err = output.RenderPage(templates, &buffer, "products.html", pageData)
		if err != nil {
			output.WriteFatalError(w, err)
		} else {
			w.WriteHeader(200)
			buffer.WriteTo(w)
		}
	}
}

func HandleGetProduct(store *Store, templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)

		product := Product{
			Title: "Dummy product",
		}

		pageData := output.PageData{
			Title: "Products",
			Data:  product,
		}

		output.RenderPage(templates, w, "product.html", pageData)
	}
}
