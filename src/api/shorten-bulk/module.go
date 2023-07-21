package api_shorten_bulk

import (
	"github.com/go-chi/chi/v5"
	api_interfaces "github.com/jeielmosi/peq.nu-backend/src/api/interfaces"
	firestore_shorten_bulk "github.com/jeielmosi/peq.nu-backend/src/core/ports/repositories/shorten-bulk/adapters/firestore"
)

type ShortenBulkModule struct {
	controller api_interfaces.Controller
	service    *ShortenBulkService
}

func (m *ShortenBulkModule) Init(r *chi.Mux) {
	m.controller.Route(r)
}

func NewShortenBulkModule() *ShortenBulkModule {
	repo := firestore_shorten_bulk.NewShortenBulkFirestore()

	service := NewShortenBulkService(&repo)
	controller := NewShortenBulkController(service)

	return &ShortenBulkModule{
		controller,
		service,
	}
}
