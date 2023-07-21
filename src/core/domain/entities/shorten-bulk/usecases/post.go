package usecases

import (
	"errors"
	"math/rand"
	"time"

	entities "github.com/jeielmosi/peq.nu-backend/src/core/domain/entities/shorten-bulk"
	helpers "github.com/jeielmosi/peq.nu-backend/src/core/helpers"
	shorten_bulk "github.com/jeielmosi/peq.nu-backend/src/core/ports/repositories/shorten-bulk/interfaces"
	types "github.com/jeielmosi/peq.nu-backend/src/core/ports/repositories/shorten-bulk/types"
)

const (
	ATTEMPTS    = 31
	OLDEST_SIZE = 101
)

func Post(
	repo *shorten_bulk.ShortenBulkRepository,
	shortenBulk entities.ShortenBulkEntity,
	hash string,
) (string, error) {
	if repo == nil {
		return "", errors.New("Repository is nil")
	}

	dto := types.NewShortenBulkRepositoryDTO(&shortenBulk)

	err := (*repo).PostSafe(hash, *dto)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func PostRandom(
	repo *shorten_bulk.ShortenBulkRepository,
	shortenBulk entities.ShortenBulkEntity,
) (string, error) {
	if repo == nil {
		return "", errors.New("Repository is nil")
	}

	dto := types.NewShortenBulkRepositoryDTO(&shortenBulk)

	for a := 0; a < ATTEMPTS; a++ {
		hash := helpers.GenerateHash()
		err := (*repo).PostSafe(hash, *dto)
		if err == nil {
			return hash, nil
		}
	}

	mpOldest, err := (*repo).GetOldest(OLDEST_SIZE)
	if err != nil {
		return "", err
	}

	mpOldestKeys := helpers.GetKeys(mpOldest)
	rnd := rand.New(
		rand.NewSource(
			time.Now().UTC().UnixNano(),
		),
	)

	for size := len(mpOldestKeys); size > 0; size-- {
		idx := rnd.Intn(size)
		hash := mpOldestKeys[idx]

		err = (*repo).PostUnsafe(hash, *dto)
		if err == nil {
			return hash, nil
		}

		mpOldestKeys[idx], mpOldestKeys[size-1] = mpOldestKeys[size-1], mpOldestKeys[idx]
	}

	return "", errors.New("Internal Error")
}
