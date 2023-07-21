package types

import (
	"errors"
	"fmt"
	"time"

	entities "github.com/jeielmosi/peq.nu-backend/src/core/domain/entities/shorten-bulk"
	repositories "github.com/jeielmosi/peq.nu-backend/src/core/ports/repositories/generics"
	helpers "github.com/jeielmosi/peq.nu-backend/src/core/ports/repositories/helpers"
	repository_helpers "github.com/jeielmosi/peq.nu-backend/src/core/ports/repositories/helpers"
)

type ShortenBulkRepositoryDTO repositories.RepositoryDTO[entities.ShortenBulkEntity]

func NewShortenBulkRepositoryDTO(shortenBulk *entities.ShortenBulkEntity) *ShortenBulkRepositoryDTO {
	now := time.Now()

	return &ShortenBulkRepositoryDTO{
		Entity:    *shortenBulk,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

const (
	CreatedAtField = "created_at"
	UpdatedAtField = "updated_at"
)

func (r *ShortenBulkRepositoryDTO) MarshalMap() (map[string]interface{}, error) {
	if r == nil {
		return nil, errors.New("ShortenBulkEntity is nil")
	}

	ans, err := r.Entity.MarshalMap()
	if err != nil {
		return nil, err
	}

	ans[CreatedAtField] = helpers.TimeToTimestamp1e8(r.CreatedAt)
	ans[UpdatedAtField] = helpers.TimeToTimestamp1e8(r.UpdatedAt)

	return ans, nil
}

func (r *ShortenBulkRepositoryDTO) UnmarshalMap(mp map[string]interface{}) error {
	if len(mp) == 0 {
		return nil
	}

	template := "The key '%s' not found on map[string]interface{}"

	createdAtStr, ok := mp[CreatedAtField].(string)
	if !ok {
		return errors.New(fmt.Sprintf(template, CreatedAtField))
	}
	createdAt, err := repository_helpers.NewTimeFromTimestamp1e8(&createdAtStr)
	if err != nil {
		return err
	}

	updatedAtStr, ok := mp[UpdatedAtField].(string)
	if !ok {
		return errors.New(fmt.Sprintf(template, UpdatedAtField))
	}
	updatedAt, err := repository_helpers.NewTimeFromTimestamp1e8(&updatedAtStr)
	if err != nil {
		return err
	}

	err = r.Entity.UnmarshalMap(mp)
	if err != nil {
		return err
	}

	r.CreatedAt = *createdAt
	r.UpdatedAt = *updatedAt

	return nil
}
