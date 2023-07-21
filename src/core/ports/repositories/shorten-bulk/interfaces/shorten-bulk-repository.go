package shorten_bulk

import (
	"time"

	types "github.com/jeielmosi/peq.nu-backend/src/core/ports/repositories/shorten-bulk/types"
)

type ShortenBulkRepository interface {
	Get(hash string) (*types.ShortenBulkRepositoryDTO, error)
	GetAndIncrement(hash string, updatedAt time.Time) (*types.ShortenBulkRepositoryDTO, error)
	GetOldest(size int) (map[string]*types.ShortenBulkRepositoryDTO, error)
	PostSafe(hash string, dto types.ShortenBulkRepositoryDTO) error
	PostUnsafe(hash string, dto types.ShortenBulkRepositoryDTO) error
}
