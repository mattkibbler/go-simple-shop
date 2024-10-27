package shop

import (
	"fmt"
	"sort"
	"sync"
)

type Store struct {
	Cache sync.Map
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) productList() ([]Product, error) {
	var products []Product
	s.Cache.Range(func(key, value any) bool {
		prod, ok := value.(Product)
		if !ok {
			return true
		}
		products = append(products, prod)
		return true
	})

	sort.SliceStable(products, func(i, j int) bool {
		return products[i].ID < products[j].ID
	})

	return products, nil
}

func (s *Store) QueryProducts(filterFunc func(p Product) bool, sortFunc func(i Product, j Product) bool) ([]Product, error) {
	products, err := s.productList()
	if err != nil {
		return []Product{}, err
	}

	var result []Product
	if filterFunc != nil {
		for _, p := range products {
			if filterFunc(p) {
				result = append(result, p)
			}
		}
	} else {
		result = products
	}

	if sortFunc != nil {
		sort.SliceStable(result, func(i, j int) bool {
			return sortFunc(result[i], result[j])
		})
	}

	return result, nil
}

func (s *Store) GetProduct(id int) (Product, error) {
	key := fmt.Sprintf("product-%d", id)
	val, ok := s.Cache.Load(key)
	if !ok {
		return Product{}, fmt.Errorf("No product found")
	}
	product, ok := val.(Product)
	if !ok {
		return Product{}, fmt.Errorf("No product found")
	}
	return product, nil

}
