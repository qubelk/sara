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

func (s *Storage) GetMany(keys ...string) map[string]string {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	tmp := make(map[string]string)

	for _, key := range keys {
		value, ok := s.store[key]
		if !ok {
			tmp[key] = "(none)"
		}

		tmp[key] = value
	}

	return tmp
}

func (s *Storage) Del(key string) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	delete(s.store, key)
}

func (s *Storage) Exists(key string) bool {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	_, ok := s.store[key]

	return ok
}
