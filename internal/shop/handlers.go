package shop

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/mattkibbler/go-simple-shop/internal/output"
)

func HandleGetProducts(store *Store, templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var buffer bytes.Buffer

		// All query params
		queryParams := r.URL.Query()
		pageStr := queryParams.Get("page")
		page, _ := strconv.Atoi(pageStr)
		if page == 0 {
			page = 1
		}
		searchQuery := queryParams.Get("search")

		products, err := store.QueryProducts(func(p Product) bool {
			if searchQuery != "" {
				return strings.Contains(strings.ToLower(p.Title), strings.ToLower(searchQuery))
			}
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
			w.WriteHeader(500)
			output.WriteFatalError(w, err)
		} else {
			w.WriteHeader(200)
			buffer.WriteTo(w)
		}
	}
}

func HandleGetProduct(store *Store, templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var buffer bytes.Buffer
		pathVars := mux.Vars(r)
		id, err := strconv.Atoi(pathVars["id"])
		if err != nil {
			output.WriteFatalError(w, err)
			return
		}

		product, err := store.GetProduct(id)
		if err != nil {
			output.WriteErrorPage(templates, w, err)
			return
		}

		pageData := output.PageData{
			Title: "Products",
			Data: struct {
				Product Product
			}{
				Product: product,
			},
		}

		err = output.RenderPage(templates, w, "product.html", pageData)
		fmt.Println(err)
		if err != nil {
			w.WriteHeader(500)
			output.WriteFatalError(w, err)
		} else {
			w.WriteHeader(200)
			buffer.WriteTo(w)
		}
	}
}
