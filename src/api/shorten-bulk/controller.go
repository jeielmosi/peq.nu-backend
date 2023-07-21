package api_shorten_bulk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	api_helpers "github.com/jeielmosi/peq.nu-backend/src/api/helpers"
	validators "github.com/jeielmosi/peq.nu-backend/src/api/validators"
)

type ShortenBulkController struct {
	service *ShortenBulkService
}

func (c *ShortenBulkController) Post(w http.ResponseWriter, r *http.Request) {
	var body validators.URLBodyDto
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		message := "Erro ao ler o corpo da requisição"
		api_helpers.ResponseJSON(&w, http.StatusBadRequest, nil, &message)
		return
	}
	validate := validators.GetValidator()
	err = validate.Struct(body)
	if err != nil {
		message := "URL inválida"
		api_helpers.ResponseJSON(&w, http.StatusBadRequest, nil, &message)
		return
	}

	params := validators.HashPostParamDto{
		Hash: chi.URLParam(r, api_helpers.HashField),
	}
	err = validate.Struct(params)
	if err != nil {
		message := "Link personalizado inválido"
		api_helpers.ResponseJSON(&w, http.StatusBadRequest, nil, &message)
		return
	}

	mp, statusCode := c.service.Post(params.Hash, body.URL)
	api_helpers.ResponseJSON(&w, statusCode, &mp, nil)
}

func (c *ShortenBulkController) PostRandom(w http.ResponseWriter, r *http.Request) {
	var body validators.URLBodyDto
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		message := "Erro ao ler o corpo da requisição"
		api_helpers.ResponseJSON(&w, http.StatusBadRequest, nil, &message)
		return
	}
	validate := validators.GetValidator()
	err = validate.Struct(body)
	if err != nil {
		message := "URL inválida"
		api_helpers.ResponseJSON(&w, http.StatusBadRequest, nil, &message)
		return
	}

	mp, statusCode := c.service.PostRandom(body.URL)
	api_helpers.ResponseJSON(&w, statusCode, &mp, nil)
}

func (c *ShortenBulkController) GetStatus(w http.ResponseWriter, r *http.Request) {
	validate := validators.GetValidator()
	params := validators.HashGetParamDto{
		Hash: chi.URLParam(r, api_helpers.HashField),
	}
	err := validate.Struct(params)
	if err != nil {
		message := "Link personalizado inválido"
		api_helpers.ResponseJSON(&w, http.StatusBadRequest, nil, &message)
		return
	}
	mp, statusCode := c.service.GetStatus(params.Hash)
	api_helpers.ResponseJSON(&w, statusCode, &mp, nil)
}

func (c *ShortenBulkController) Get(w http.ResponseWriter, r *http.Request) {
	validate := validators.GetValidator()
	params := validators.HashGetParamDto{
		Hash: chi.URLParam(r, api_helpers.HashField),
	}
	err := validate.Struct(params)
	if err != nil {
		message := "Link personalizado inválido"
		api_helpers.ResponseJSON(&w, http.StatusBadRequest, nil, &message)
		return
	}
	mp, statusCode := c.service.Get(params.Hash)
	api_helpers.ResponseJSON(&w, statusCode, &mp, nil)
}

func (c *ShortenBulkController) Route(r *chi.Mux) {

	r.Get(fmt.Sprintf("/{%s}/status", api_helpers.HashField), c.GetStatus)
	r.Get(fmt.Sprintf("/{%s}", api_helpers.HashField), c.Get)

	r.Group(func(r chi.Router) {
		r.Use(httprate.LimitByIP(10, time.Minute))
		r.Post("/", c.PostRandom)
		r.Post(fmt.Sprintf("/{%s}", api_helpers.HashField), c.Post)
	})
}

func NewShortenBulkController(service *ShortenBulkService) *ShortenBulkController {
	return &ShortenBulkController{
		service,
	}
}
