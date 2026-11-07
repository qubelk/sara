package storage

import (
	"sync"
)

type Storage struct {
	store map[string]string
	mtx   sync.RWMutex
}

func New() *Storage {
	return &Storage{
		store: make(map[string]string),
		mtx:   sync.RWMutex{},
	}
}

func (s *Storage) Set(key, value string) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.store[key] = value
}

func (s *Storage) Get(key string) string {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	value, ok := s.store[key]
	if !ok {
		return ""
	}

	return value
}
