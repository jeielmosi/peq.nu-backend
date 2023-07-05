package api_shorten_bulk

import (
	"net/http"

	api_helpers "github.com/jei-el/vuo.be-backend/src/api/helpers"
	entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	usecases "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk/usecases"
	shorten_bulk_gateway "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/gateways/interfaces"
)

type ShortenBulkService struct {
	gateway *shorten_bulk_gateway.ShortenBulkGateway
}

func (s *ShortenBulkService) Post(url string) (map[string]interface{}, int) {
	ans := map[string]interface{}{}

	if url == "" {
		return ans, http.StatusBadRequest
	}

	shortenBulk := entities.NewShortenBulkEntity(url, 0)

	hash, err := usecases.Post(s.gateway, *shortenBulk)
	if err != nil {
		ans[api_helpers.MessageField] = err.Error()
		return ans, http.StatusInternalServerError
	}

	ans[api_helpers.HashField] = hash
	return ans, http.StatusCreated
}

func (s *ShortenBulkService) Get(hash string) (map[string]interface{}, int) {
	ans := map[string]interface{}{}

	if hash == "" {
		return ans, http.StatusNotFound
	}

	shortenBulk, err := usecases.Get(s.gateway, hash)
	if err != nil {
		ans[api_helpers.MessageField] = err.Error()
		return ans, http.StatusInternalServerError
	}

	return usecases.ToMapInterface(shortenBulk), http.StatusOK
}

func NewShortenBulkService(gateway *shorten_bulk_gateway.ShortenBulkGateway) *ShortenBulkService {
	return &ShortenBulkService{
		gateway,
	}
}
