package api_shorten_bulk

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	api_helpers "github.com/jei-el/vuo.be-backend/src/api/helpers"
)

type ShortenBulkController struct {
	service *ShortenBulkService
}

func (c *ShortenBulkController) Post(w http.ResponseWriter, r *http.Request) {
	body := map[string]interface{}{}
	json.NewDecoder(r.Body).Decode(&body)

	iUrl := body[api_helpers.URLField]

	url, ok := iUrl.(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatalf("Field %s not found", api_helpers.URLField)
	}

	mp, statusCode := c.service.Post(url)

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	res, err := json.Marshal(mp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("Error happened in JSON marshal. Err: %s", err.Error())
	}

	w.Write(res)
}

func (c *ShortenBulkController) Get(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, api_helpers.HashField)
	mp, statusCode := c.service.Get(hash)

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	res, err := json.Marshal(mp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err.Error())
	}

	w.Write(res)
}

func (c *ShortenBulkController) Route(r *chi.Mux) {

	r.Get(fmt.Sprintf("/{%s}", api_helpers.HashField), c.Get)
	r.Post("/", c.Post)
}

func NewShortenBulkController(service *ShortenBulkService) *ShortenBulkController {
	return &ShortenBulkController{
		service,
	}
}
