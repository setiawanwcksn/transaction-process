package api

import (
	"flip/internal/repository/bus"
	"flip/internal/repository/csvparser"
	"flip/internal/repository/storage"
	"net/http"
)

type Handler struct {
	storage storage.Repository
	parser  csvparser.Repository
	bus     bus.Repository
	async   bool
}

func New(storage storage.Repository, parser csvparser.Repository, bus bus.Repository) *Handler {
	return &Handler{
		storage: storage,
		parser:  parser,
		bus:     bus,
		async:   true,
	}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/health", h.health)
	mux.HandleFunc("/statements", h.upload)
	mux.HandleFunc("/balance", h.balance)
	mux.HandleFunc("/transactions/issues", h.issues)
	mux.HandleFunc("/progress", h.progress)
}
