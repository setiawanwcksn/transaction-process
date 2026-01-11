package worker

import (
	"context"
	"flip/internal/model"
	"sync"
	"time"
)

type HandleResult struct {
	Success   bool
	Retryable bool
}

type Repository struct {
	bus     <-chan model.Event
	tracker *tracker
	handler func(context.Context, model.Event) *HandleResult
	stop    chan struct{}
	sleep   func(time.Duration)
}

type tracker struct {
	mu   sync.Mutex
	seen map[string]bool
}

func newTracker() *tracker {
	return &tracker{seen: make(map[string]bool)}
}

func StartDefault(bus <-chan model.Event, handler func(context.Context, model.Event) *HandleResult) *Repository {
	return New(bus, handler)
}

func New(bus <-chan model.Event, handler func(context.Context, model.Event) *HandleResult) *Repository {
	return &Repository{
		bus:     bus,
		tracker: newTracker(),
		handler: handler,
		stop:    make(chan struct{}),
		sleep:   time.Sleep,
	}
}
