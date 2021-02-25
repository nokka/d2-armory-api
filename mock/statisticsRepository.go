package mock

import (
	"context"

	"github.com/nokka/d2-armory-api/internal/domain"
)

// StatisticsRepository is a mock implementation of the statistics repository.
type StatisticsRepository struct {
	GetByCharacterFn      func(ctx context.Context, character string) (*domain.CharacterStatistics, error)
	GetByCharacterInvoked bool
	UpsertFn              func(ctx context.Context, stat domain.StatisticsRequest) error
	UpsertInvoked         bool
}

// GetByCharacter calls the GetByCharacterFn and registers the invoke.
func (r *StatisticsRepository) GetByCharacter(ctx context.Context, character string) (*domain.CharacterStatistics, error) {
	r.GetByCharacterInvoked = true
	return r.GetByCharacterFn(ctx, character)
}

// Upsert calls the UpsertFn and registers the invoke.
func (r *StatisticsRepository) Upsert(ctx context.Context, stat domain.StatisticsRequest) error {
	r.UpsertInvoked = true
	return r.UpsertFn(ctx, stat)
}
