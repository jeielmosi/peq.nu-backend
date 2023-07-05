package shorten_bulk_gateway

import (
	entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
)

type ShortenBulkGateway interface {
	Get(hash string) (*entities.ShortenBulkEntity, error)
	Post(shorten_bulk entities.ShortenBulkEntity) (string, error)
}
