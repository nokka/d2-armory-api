// Package character contains the heart of the domain model.
package character

import (
	"errors"
	"time"

	"github.com/nokka/d2s"
)

// Character is the central class in the domain model.
type Character struct {
	ID         string         `json:"d2s_id"`
	D2s        *d2s.Character `json:"d2s"`
	LastParsed time.Time      `json:"last_parsed"`
}

// New creates a new, unparsed Character.
func New(id string, char d2s.Character) *Character {
	return &Character{
		ID:         id,
		D2s:        &char,
		LastParsed: time.Now(),
	}
}

// Repository provides access a character store.
type Repository interface {
	Find(string) *Character
	Update(character *Character) error
	Store(character *Character) error
}

// ErrUnknown is used when a character could not be found.
var ErrUnknown = errors.New("Unknown character")
