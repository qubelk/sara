package storage

import "sync"

type Storage struct {
	store map[string]string
	mtx   sync.Mutex
}

func New() *Storage {
	return &Storage{
		store: make(map[string]string),
		mtx:   sync.Mutex{},
	}
}
