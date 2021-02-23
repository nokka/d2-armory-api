package statistics

import (
	"context"

	"github.com/nokka/d2-armory-api/internal/domain"
)

// characterRepository is the interface representation of the data layer
// the service depend on.
type statisticsRepository interface {
	GetByCharacter(ctx context.Context, character string) (*domain.StatisticsRequest, error)
	Upsert(ctx context.Context, stat domain.StatisticsRequest) error
}

// Service performs all operations on statistics.
type Service struct {
	repository statisticsRepository
}

// GetCharacter will get the statistics on a specific character.
func (s Service) GetCharacter(ctx context.Context, character string) (*domain.StatisticsRequest, error) {
	char, err := s.repository.GetByCharacter(ctx, character)
	if err != nil {
		return nil, err
	}

	return char, nil
}

// Parse will perform the storage of a statistics request.
func (s Service) Parse(ctx context.Context, stats []domain.StatisticsRequest) error {
	for _, req := range stats {
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
