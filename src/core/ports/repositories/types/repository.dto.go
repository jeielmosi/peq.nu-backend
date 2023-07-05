package repositories

import (
	"time"
)

type RepositoryDTO[T any] struct {
	Entity    *T
	Locked    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r *RepositoryDTO[T]) Update() *RepositoryDTO[T] {
	now := time.Now()
	return &RepositoryDTO[T]{
		Entity:    r.Entity,
		CreatedAt: r.CreatedAt,
		Locked:    r.Locked,
		UpdatedAt: now,
	}
}

func (r *RepositoryDTO[T]) LockSwitch() *RepositoryDTO[T] {
	now := time.Now()
	return &RepositoryDTO[T]{
		Entity:    r.Entity,
		CreatedAt: r.CreatedAt,
		Locked:    !r.Locked,
		UpdatedAt: now,
	}
}

func NewRepositoryDTO[T any](
	entity *T,
	locked bool,
) *RepositoryDTO[T] {
	now := time.Now()

	return &RepositoryDTO[T]{
		Entity:    entity,
		CreatedAt: now,
		Locked:    locked,
		UpdatedAt: now,
	}
}
