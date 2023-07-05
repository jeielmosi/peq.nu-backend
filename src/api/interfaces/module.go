package api_interfaces

import (
	"github.com/go-chi/chi/v5"
)

type Module interface {
	Init(r *chi.Mux)
}
