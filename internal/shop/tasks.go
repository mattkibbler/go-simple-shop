package shop

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

func SerializeProductCache(filePath string, productCache *sync.Map) error {
	tempMap := make(map[string]Product)

	productCache.Range(func(key, value interface{}) bool {
		k, ok := key.(string)
		if !ok {
			return false
		}
		v, ok := value.(Product)
		if !ok {
			return false
		}
		tempMap[k] = v
		return true
	})

	data, err := json.Marshal(tempMap)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func UnserializeProductCache(productCachePath string, productCache *sync.Map) (int, error) {
	// check file exists, if not then this is probably the first time the server has ran
	if _, err := os.Stat(productCachePath); os.IsNotExist(err) {
		return 0, nil
	}

	// read file data
	data, err := os.ReadFile(productCachePath)
	if err != nil {
		return 0, fmt.Errorf("failed to read product cache file: %v", err)
	}

	// deserialize JSON into an ordinary map
	tempMap := make(map[string]Product)
	if err := json.Unmarshal(data, &tempMap); err != nil {
		return 0, fmt.Errorf("failed to deserialize product cache: %v", err)
	}

	// restore data into the original sync.map
	pCount := 0
	for key, value := range tempMap {
		productCache.Store(key, value)
		pCount++
	}
	return pCount, nil
}
