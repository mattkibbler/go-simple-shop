package shop

import (
	"sync"
)

type Store struct {
	Cache sync.Map
}

func NewStore() *Store {
	return &Store{}
}
