package statistics

import (
	"context"
	"fmt"
	"strings"

	"github.com/nokka/d2-armory-api/internal/domain"
)

var validDifficulties = map[string]struct{}{
	domain.DifficultyNormal:    {},
	domain.DifficultyNightmare: {},
	domain.DifficultyHell:      {},
}

// characterRepository is the interface representation of the data layer
// the service depend on.
type statisticsRepository interface {
	GetByCharacter(ctx context.Context, character string) (*domain.CharacterStatistics, error)
	Upsert(ctx context.Context, stat domain.StatisticsRequest) error
}

// Service performs all operations on statistics.
type Service struct {
	repository statisticsRepository
}

// GetCharacter will get the statistics on a specific character.
func (s Service) GetCharacter(ctx context.Context, character string) (*domain.CharacterStatistics, error) {
	char, err := s.repository.GetByCharacter(ctx, character)
	if err != nil {
		return nil, err
	}

	return char, nil
}

// Parse will perform the storage of a statistics request.
func (s Service) Parse(ctx context.Context, stats []domain.StatisticsRequest) error {
	for _, req := range stats {
		if _, valid := validDifficulties[req.Difficulty]; !valid {
			return fmt.Errorf("difficulty %s supplied for character %s, %w", req.Difficulty, req.Character, domain.ErrRequest)
		}

		// Lower case the character name to keep consistency.
		req.Account = strings.ToLower(req.Account)
		req.Character = strings.ToLower(req.Character)

		// Upsert each character stat request.
		err := s.repository.Upsert(ctx, req)
		if err != nil {
			return err
		}
	}

	return nil
}

// NewService constructs a new statistics service with all the dependencies.
func NewService(repository statisticsRepository) *Service {
	return &Service{
		repository: repository,
	}
}
