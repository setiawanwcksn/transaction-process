package csvparser

import (
	"flip/internal/repository/bus"
	"flip/internal/repository/storage"
	"flip/internal/repository/worker"
)

type Repository struct {
	storage storage.Repository
	bus     bus.Repository
	worker  worker.Repository
}

func NewParser(storage storage.Repository, bus bus.Repository, worker worker.Repository) *Repository {
	return &Repository{
		storage: storage,
		bus:     bus,
		worker:  worker,
	}
}
