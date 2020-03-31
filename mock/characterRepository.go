package mock

import (
	"context"

	"github.com/nokka/d2-armory-api/internal/domain"
)

// CharacterRepository is a mock implementation of the character repository.
type CharacterRepository struct {
	FindFn        func(ctx context.Context, id string) (*domain.Character, error)
	FindInvoked   bool
	UpdateFn      func(ctx context.Context, character *domain.Character) error
	UpdateInvoked bool
	StoreFn       func(ctx context.Context, character *domain.Character) error
	StoreInvoked  bool
}

// Find calls the FindFn and registers the invoke.
func (r CharacterRepository) Find(ctx context.Context, id string) (*domain.Character, error) {
	r.FindInvoked = true
	return r.FindFn(ctx, id)
}

// Update calls the UpdateFn and registers the invoke.
func (r CharacterRepository) Update(ctx context.Context, character *domain.Character) error {
	r.UpdateInvoked = true
	return r.UpdateFn(ctx, character)
}

// Store calls the StoreFn and registers the invoke.
func (r CharacterRepository) Store(ctx context.Context, character *domain.Character) error {
	r.StoreInvoked = true
	return r.StoreFn(ctx, character)
}
