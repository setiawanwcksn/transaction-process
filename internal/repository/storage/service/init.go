package storage

import (
	"sync"

	"flip/internal/model"
)

type Repository struct {
	mu   sync.RWMutex
	data map[string]*model.Statement
}

func NewStorage() *Repository {
	return &Repository{
		data: make(map[string]*model.Statement),
	}
}
