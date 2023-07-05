package api_interfaces

import "github.com/go-chi/chi/v5"

type Controller interface {
	Route(r *chi.Mux)
}
