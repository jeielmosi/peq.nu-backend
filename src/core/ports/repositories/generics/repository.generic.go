package repositories

import (
	"time"
)

type RepositoryDTO[T any] struct {
	Entity    T
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r *RepositoryDTO[T]) Update() *RepositoryDTO[T] {
	now := time.Now()
	return &RepositoryDTO[T]{
		Entity:    r.Entity,
		CreatedAt: r.CreatedAt,
		UpdatedAt: now,
	}
}

func NewRepositoryDTO[T any](entity *T) *RepositoryDTO[T] {
	now := time.Now()

	return &RepositoryDTO[T]{
		Entity:    *entity,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
