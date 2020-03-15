package character

import (
	"fmt"
	"regexp"

	"github.com/nokka/d2-armory-api/internal/domain"
)

// Dependencies.
type characterRepository interface {
	Find(id string) (*domain.Character, error)
	Update(character *domain.Character) error
	Store(character *domain.Character) error
}

// Service performs all operations on parsing characters.
type Service struct {
	d2spath    string
	characters characterRepository
}

// The name regexp required for character names, to enforce strict diablo rules
// on the names to prevent missuse of the endpoint.
const nameRegexp = "^[a-zA-Z]+[_-]?[a-zA-Z]+$"

// Parse will perform the actual parsing of the character.
func (s Service) Parse(name string) (*domain.Character, error) {
	match, _ := regexp.MatchString(nameRegexp, name)
	if !match {
		return nil, domain.ErrInvalidArgument
	}

	c, err := s.characters.Find(name)
	if err != nil {
		return nil, err
	}

	fmt.Println("AFTER FIND")
	fmt.Println(c)

	return &domain.Character{}, nil
}

// NewService constructs a new parsing service with all the dependencies.
func NewService(d2spath string, characterRepository characterRepository) *Service {
	return &Service{
		d2spath:    d2spath,
		characters: characterRepository,
	}
}
