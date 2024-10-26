package shop

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func FetchAndCacheProducts(productCache *sync.Map) error {
	prodLimit := 30
	totalProducts := 999 // Set initial high value
	fetchedProducts := 0
	for fetchedProducts < totalProducts {
		url := fmt.Sprintf("https://dummyjson.com/products?limit=%d&skip=%d", prodLimit, fetchedProducts)
		fmt.Println(url)
		resp, err := http.Get(url)
		if err != nil {
			return err
		}

		var productsResponse ProductAPIResponse
		if err := json.NewDecoder(resp.Body).Decode(&productsResponse); err != nil {
			return fmt.Errorf("could not decode product JSON: %v", err)
		}

		for _, product := range productsResponse.Products {
			productCache.Store(fmt.Sprintf("product-%d", product.ID), product)
		}

		totalProducts = productsResponse.Total
		fetchedProducts += len(productsResponse.Products)

		// Close response body
		resp.Body.Close()

		// Wait for a moment so we're not overloading the API...
		time.Sleep(time.Second)
	}

	return nil
}
