package usecases

import (
	"errors"

	entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	shorten_bulk_gateway "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/gateways/interfaces"
)

func Post(
	gateway *shorten_bulk_gateway.ShortenBulkGateway,
	shortenBulk entities.ShortenBulkEntity,
) (string, error) {
	if gateway == nil {
		return "", errors.New("Gateway is nil")
	}
	return (*gateway).Post(shortenBulk)
}
