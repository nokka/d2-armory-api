package character

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/nokka/d2-armory-api/internal/domain"
)

//go:generate moq -out ./service_mocks.go . parser characterRepository

// parser is the interface representation of a d2 parser the service depend on.
type parser interface {
	Parse(name string) (*domain.Character, error)
}

// characterRepository is the interface representation of the data layer
// the service depend on.
type characterRepository interface {
	Find(ctx context.Context, id string) (*domain.Character, error)
	Update(ctx context.Context, character *domain.Character) error
	Store(ctx context.Context, character *domain.Character) error
}

// Service performs all operations on parsing characters.
type Service struct {
	parser        parser
	characters    characterRepository
	cacheDuration time.Duration
}

// The name regexp required for character names, to enforce strict diablo rules
// on the names to prevent missuse of the endpoint.
const nameRegexp = "^[a-zA-Z]+[_-]?[a-zA-Z]+$"

// Parse will perform the actual parsing of the character.
func (s Service) Parse(ctx context.Context, name string) (*domain.Character, error) {
	match, _ := regexp.MatchString(nameRegexp, name)
	if !match {
		return nil, domain.ErrInvalidArgument
	}

	// Read character from db cache.
	c, err := s.characters.Find(ctx, name)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			// Character didn't exist at all, so lets parse and store it.
			parsed, err := s.parser.Parse(name)
			if err != nil {
				return nil, err
			}

			if err := s.characters.Store(ctx, parsed); err != nil {
				return nil, err
			}

			return parsed, nil
		}

		// The error wasn't ErrNotFound, so just return it.
		return nil, err
	}

	// Character already exists, let's check how long since we parsed it.
	diff := time.Since(c.LastParsed)

	if diff >= s.cacheDuration {
		parsed, err := s.parser.Parse(name)
		if err != nil {
			return nil, err
		}

		// Update the existing record in the db.
		err = s.characters.Update(ctx, parsed)
		if err != nil {
			return nil, err
		}

		return parsed, nil
	}

	// We parsed this character less than 3 minutes ago so return the db version.
	return c, nil
}

// NewService constructs a new parsing service with all the dependencies.
func NewService(parser parser, characterRepository characterRepository, cacheDuration time.Duration) *Service {
	return &Service{
		parser:        parser,
		characters:    characterRepository,
		cacheDuration: cacheDuration,
	}
}
