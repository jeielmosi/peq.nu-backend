package api_shorten_bulk

import (
	"net/http"
	"time"

	api_helpers "github.com/jeielmosi/peq.nu-backend/src/api/helpers"
	entities "github.com/jeielmosi/peq.nu-backend/src/core/domain/entities/shorten-bulk"
	usecases "github.com/jeielmosi/peq.nu-backend/src/core/domain/entities/shorten-bulk/usecases"
	shorten_bulk "github.com/jeielmosi/peq.nu-backend/src/core/ports/repositories/shorten-bulk/interfaces"
)

type ShortenBulkService struct {
	repo *shorten_bulk.ShortenBulkRepository
}

func (s *ShortenBulkService) Post(hash string, url string) (map[string]interface{}, int) {
	ans := make(map[string]interface{})
	shortenBulk := entities.NewShortenBulkEntity(url, 0, true)

	hash, err := usecases.Post(s.repo, *shortenBulk, hash)
	if err != nil {
		ans[api_helpers.MessageField] = err.Error()
		return ans, http.StatusInternalServerError
	}

	ans[api_helpers.HashField] = hash
	return ans, http.StatusCreated
}

func (s *ShortenBulkService) PostRandom(url string) (map[string]interface{}, int) {
	ans := make(map[string]interface{})
	shortenBulk := entities.NewShortenBulkEntity(url, 0, false)

	hash, err := usecases.PostRandom(s.repo, *shortenBulk)
	if err != nil {
		ans[api_helpers.MessageField] = err.Error()
		return ans, http.StatusInternalServerError
	}

	ans[api_helpers.HashField] = hash
	return ans, http.StatusCreated
}

func (s *ShortenBulkService) GetStatus(hash string) (map[string]interface{}, int) {
	ans := make(map[string]interface{})

	if hash == "" {
		return ans, http.StatusNotFound
	}

	shortenBulk, err := usecases.GetStatus(s.repo, hash)
	if err != nil {
		ans[api_helpers.MessageField] = err.Error()
		return ans, http.StatusInternalServerError
	}

	flatten, err := shortenBulk.MarshalMap()
	if err != nil {
		ans[api_helpers.MessageField] = err.Error()
		return ans, http.StatusInternalServerError
	}

	return flatten, http.StatusOK
}

func (s *ShortenBulkService) Get(hash string) (map[string]interface{}, int) {
	now := time.Now()
	ans := make(map[string]interface{})

	if hash == "" {
		return ans, http.StatusNotFound
	}

	shortenBulk, err := usecases.Get(s.repo, hash, now)
	if err != nil {
		ans[api_helpers.MessageField] = err.Error()
		return ans, http.StatusInternalServerError
	}

	flatten, err := shortenBulk.MarshalMap()
	if err != nil {
		ans[api_helpers.MessageField] = err.Error()
		return ans, http.StatusInternalServerError
	}

	return flatten, http.StatusOK
}

func NewShortenBulkService(repo *shorten_bulk.ShortenBulkRepository) *ShortenBulkService {
	return &ShortenBulkService{
		repo,
	}
}
