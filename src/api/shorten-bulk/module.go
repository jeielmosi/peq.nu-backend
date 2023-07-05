package api_shorten_bulk

import (
	"log"
	"os"

	"github.com/go-chi/chi/v5"
	api_interfaces "github.com/jei-el/vuo.be-backend/src/api/interfaces"
	config "github.com/jei-el/vuo.be-backend/src/config"
	GAM "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/gateways/GAM"
)

type ShortenBulkModule struct {
	controller api_interfaces.Controller
	service    *ShortenBulkService
}

func (m *ShortenBulkModule) Init(r *chi.Mux) {
	m.controller.Route(r)
}

func NewShortenBulkModule() *ShortenBulkModule {
	envName := os.Getenv(config.CURRENT_ENV)
	gateway, err := GAM.NewGAMShortenBulkGateway(envName)
	if err != nil {
		log.Println("Error on create a gatway on ShortenBulkModule")
		log.Fatalln(err.Error())
	}

	service := NewShortenBulkService(&gateway)
	controller := NewShortenBulkController(service)

	return &ShortenBulkModule{
		controller,
		service,
	}
}
