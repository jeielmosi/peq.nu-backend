package usecases

import (
	"errors"
	"time"

	entities "github.com/jeielmosi/peq.nu-backend/src/core/domain/entities/shorten-bulk"
	shorten_bulk "github.com/jeielmosi/peq.nu-backend/src/core/ports/repositories/shorten-bulk/interfaces"
)

func Get(repo *shorten_bulk.ShortenBulkRepository, hash string, now time.Time) (
	*entities.ShortenBulkEntity,
	error,
) {
	if repo == nil {
		return nil, errors.New("Repository is nil")
	}
	dto, err := (*repo).GetAndIncrement(hash, now)
	if err != nil {
		return nil, err
	}

	return &dto.Entity, nil
}

func GetStatus(repo *shorten_bulk.ShortenBulkRepository, hash string) (
	*entities.ShortenBulkEntity,
	error,
) {
	if repo == nil {
		return nil, errors.New("Repository is nil")
	}
	dto, err := (*repo).Get(hash)
	if err != nil {
		return nil, err
	}

	return &dto.Entity, nil
}
