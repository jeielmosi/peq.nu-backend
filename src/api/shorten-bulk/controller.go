package api_shorten_bulk

import (
	"encoding/json"
	"fmt"
	"log"
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	validate := validators.GetValidator()
	err = validate.Struct(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := validators.HashPostParamDto{
		Hash: chi.URLParam(r, api_helpers.HashField),
	}
	err = validate.Struct(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mp, statusCode := c.service.Post(params.Hash, body.URL)

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	res, err := json.Marshal(mp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("Error happened in JSON marshal. Err: %s\n", err.Error())
	}

	w.Write(res)
}

func (c *ShortenBulkController) PostRandom(w http.ResponseWriter, r *http.Request) {
	var body validators.URLBodyDto
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	validate := validators.GetValidator()
	err = validate.Struct(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mp, statusCode := c.service.PostRandom(body.URL)

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	res, err := json.Marshal(mp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("Error happened in JSON marshal. Err: %s\n", err.Error())
	}

	w.Write(res)
}

func (c *ShortenBulkController) GetStatus(w http.ResponseWriter, r *http.Request) {
	validate := validators.GetValidator()
	params := validators.HashGetParamDto{
		Hash: chi.URLParam(r, api_helpers.HashField),
	}
	err := validate.Struct(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mp, statusCode := c.service.GetStatus(params.Hash)

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	res, err := json.Marshal(mp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s\n", err.Error())
	}

	w.Write(res)
}

func (c *ShortenBulkController) Get(w http.ResponseWriter, r *http.Request) {
	validate := validators.GetValidator()
	params := validators.HashGetParamDto{
		Hash: chi.URLParam(r, api_helpers.HashField),
	}
	err := validate.Struct(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mp, statusCode := c.service.Get(params.Hash)

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	res, err := json.Marshal(mp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s\n", err.Error())
	}

	w.Write(res)
}

func (c *ShortenBulkController) Route(r *chi.Mux) {

	r.Get(fmt.Sprintf("/{%s}/status", api_helpers.HashField), c.GetStatus)
	r.Get(fmt.Sprintf("/{%s}", api_helpers.HashField), c.Get)

	r.Group(func(r chi.Router) {
		r.Use(httprate.LimitByIP(10, time.Hour))
		r.Post("/", c.PostRandom)
		r.Post(fmt.Sprintf("/{%s}", api_helpers.HashField), c.Post)
	})
}

func NewShortenBulkController(service *ShortenBulkService) *ShortenBulkController {
	return &ShortenBulkController{
		service,
	}
}
