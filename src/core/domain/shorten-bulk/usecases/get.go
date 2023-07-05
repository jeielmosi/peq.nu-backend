package usecases

import (
	"errors"

	entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	shorten_bulk_gateway "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/gateways/interfaces"
)

func Get(gateway *shorten_bulk_gateway.ShortenBulkGateway, hash string) (
	*entities.ShortenBulkEntity,
	error,
) {
	if gateway == nil {
		return nil, errors.New("Gateway is nil")
	}
	return (*gateway).Get(hash)
}
