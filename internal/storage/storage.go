package storage

import (
	"maps"
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

func (s *Storage) Get(keys ...string) map[string]string {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	tmp := make(map[string]string)

	for _, key := range keys {
		tmp[key] = s.store[key]
	}

	return tmp
}

func (s *Storage) Del(keys ...string) map[string]bool {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	tmp := make(map[string]bool)

	for _, key := range keys {
		_, ok := s.store[key]
		delete(s.store, key)
		tmp[key] = ok
	}

	return tmp
}

func (s *Storage) Exists(keys ...string) map[string]bool {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	tmp := make(map[string]bool)

	for _, key := range keys {
		_, ok := s.store[key]
		tmp[key] = ok
	}

	return tmp
}

func (s *Storage) Dump() map[string]string {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	tmp := make(map[string]string)
	maps.Copy(tmp, s.store)

	return tmp
}
