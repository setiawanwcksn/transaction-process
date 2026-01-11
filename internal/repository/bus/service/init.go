package bus

import "flip/internal/model"

type Repository struct {
	ch chan model.Event
}

func New(buffer int) *Repository {
	return &Repository{
		ch: make(chan model.Event, buffer),
	}
}

func NewDefault() *Repository {
	return New(100)
}