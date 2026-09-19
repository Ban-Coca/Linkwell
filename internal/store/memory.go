package store

import (
	"sync"
	
)

type MemoryStore struct {
	mu sync.RWMutex
	links map[string]Link
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		links: make(map[string]Link),
	}
}

func (s *MemoryStore) Create(link Link) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.links[link.Code]; ok {
		return ErrAlreadyExists
	}

	s.links[link.Code] = link
	return nil
}

func (s *MemoryStore) Get(code string) (Link, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	link, ok := s.links[code]
	if !ok {
		return Link{}, ErrNotFound
	}
	return link, nil
}

func (s *MemoryStore) Exists(code string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.links[code]
	return ok, nil
}

var _ LinkStore = (*MemoryStore)(nil)