package funcs

import (
	"errors"
	"time"

	entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	shorten_bulk "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/interfaces"
	repositories "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/types"
)

func NewUnlockFunc(hash string, updatedAt time.Time) func(*shorten_bulk.ShortenBulkRepository) (
	*repositories.RepositoryDTO[entities.ShortenBulkEntity],
	error,
) {
	return func(repository *shorten_bulk.ShortenBulkRepository) (
		*repositories.RepositoryDTO[entities.ShortenBulkEntity],
		error,
	) {
		if repository == nil {
			return nil, errors.New("Repository is nil")
		}
		err := (*repository).Unlock(hash, updatedAt)
		return nil, err
	}
}
