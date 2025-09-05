package repository

import (
	"sync"

	"github.com/RoGogDBD/kasp/internal/models"
)

type Storage struct {
	mu     sync.RWMutex
	status map[string]models.Status
}

func NewStorage() *Storage {
	return &Storage{status: make(map[string]models.Status)}
}

func (s *Storage) SetStatus(id string, st models.Status) {
	s.mu.Lock()
	s.status[id] = st
	s.mu.Unlock()
}

func (s *Storage) GetStatus(id string) (models.Status, bool) {
	s.mu.RLock()
	st, ok := s.status[id]
	s.mu.RUnlock()
	return st, ok
}
