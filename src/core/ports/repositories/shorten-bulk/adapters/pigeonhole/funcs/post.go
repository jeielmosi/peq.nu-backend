package funcs

import (
	"errors"

	entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	shorten_bulk "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/interfaces"
	repositories "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/types"
)

func NewPostFunc(hash string, dto repositories.RepositoryDTO[entities.ShortenBulkEntity]) func(*shorten_bulk.ShortenBulkRepository) (
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
		err := (*repository).Post(hash, dto)
		return nil, err
	}
}
