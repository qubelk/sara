package storage

import (
	"maps"
	"sync"
	"time"
)

type item struct {
	Value      string
	ExpireTime time.Time
}

type Storage struct {
	store map[string]item
	mtx   sync.RWMutex
}

func New() *Storage {
	s := &Storage{
		store: make(map[string]item),
		mtx:   sync.RWMutex{},
	}

	go s.janitor(200 * time.Millisecond)

	return s
}

func (s *Storage) janitor(interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()

	for range t.C {
		s.CleanAllExpired()
	}
}

func (s *Storage) Set(key, value string, ttl time.Duration) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	var expireTime time.Time
	if ttl > 0 {
		expireTime = time.Now().Add(ttl)
	}

	s.store[key] = item{
		Value:      value,
		ExpireTime: expireTime,
	}
}

func (s *Storage) Get(keys ...string) map[string]item {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	tmp := make(map[string]item)

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

func (s *Storage) Dump() map[string]item {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	tmp := make(map[string]item)
	maps.Copy(tmp, s.store)

	return tmp
}

func (s *Storage) CleanAllExpired() {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	now := time.Now()

	for k, v := range s.store {
		if !v.ExpireTime.IsZero() && now.After(v.ExpireTime) {
			delete(s.store, k)
		}
	}
}
