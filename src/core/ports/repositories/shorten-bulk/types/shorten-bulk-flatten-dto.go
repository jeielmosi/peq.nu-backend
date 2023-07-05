package types

import (
	entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	helpers "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/helpers"
	repository_helpers "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/helpers"
	repositories "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/types"
)

const (
	CreatedAtField = "created_at"
	UpdatedAtField = "updated_at"
	URLField       = "url"
	LockedField    = "locked"
	ClicksField    = "clicks"
)

type ShortenBulkFlattenDTO = map[string]interface{}

func NewShortenBulkFlattenDTO(
	dto repositories.RepositoryDTO[entities.ShortenBulkEntity],
) ShortenBulkFlattenDTO {
	ans := ShortenBulkFlattenDTO{}

	ans[CreatedAtField] = helpers.TimeToTimestamp1e8(dto.CreatedAt)
	ans[UpdatedAtField] = helpers.TimeToTimestamp1e8(dto.UpdatedAt)
	ans[URLField] = dto.Entity.URL
	ans[LockedField] = dto.Locked
	ans[ClicksField] = dto.Entity.Clicks

	return ans
}

func ToRepositoryDTO(flatten ShortenBulkFlattenDTO) (
	*repositories.RepositoryDTO[entities.ShortenBulkEntity],
	error,
) {
	if len(flatten) == 0 {
		return nil, nil
	}

	createdAtStr, ok := flatten[CreatedAtField].(string)
	if !ok {
		return nil, nil
	}

	createdAt, err := repository_helpers.NewTimeFromTimestamp1e8(&createdAtStr)
	if err != nil {
		return nil, err
	}

	updatedAtStr, ok := flatten[UpdatedAtField].(string)
	if !ok {
		return nil, nil
	}
	updatedAt, err := repository_helpers.NewTimeFromTimestamp1e8(&updatedAtStr)
	if err != nil {
		return nil, err
	}

	dto := &repositories.RepositoryDTO[entities.ShortenBulkEntity]{
		Entity: entities.NewShortenBulkEntity(
			flatten[URLField].(string),
			flatten[ClicksField].(int64),
		),
		CreatedAt: *createdAt,
		Locked:    flatten[LockedField].(bool),
		UpdatedAt: *updatedAt,
	}

	return dto, err
}
