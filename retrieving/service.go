// Package retrieving is responsible for retrieving and parsing character data
// from the diablo character d2s files.
package retrieving

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nokka/armory/character"
	"github.com/nokka/d2s"
)

// ErrInvalidArgument is returned when one or more arguments are invalid.
var ErrInvalidArgument = errors.New("invalid argument")

// ErrNonExistingCharacter is returned when someone requests data for a d2s char
// that does not exists on disk.
var ErrNonExistingCharacter = errors.New("The requested character does not exist")

// Service provides operations on d2s character data.
type Service interface {
	// GetCharacter will return the character with the given name.
	RetrieveCharacter(string) (character.Character, error)
}

type service struct {
	characters character.Repository
}

func (s *service) RetrieveCharacter(name string) (character.Character, error) {
	if name == "" {
		return character.Character{}, ErrInvalidArgument
	}

	path := "/Users/stekon/go/src/github.com/nokka/armory/testdata/0bb3677cf40d3adc.nokkazon"
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Error while opening .d2s file", err)
	}

	defer file.Close()

	char, err := d2s.Parse(file)

	if err != nil {
		fmt.Println(err)
	}

	/*var result []api.Character
	for _, c := range s.characters.FindByType(ladderType) {
		result = append(result, *c)
	}*/

	return character.Character{
		ID:         "test",
		D2s:        &char,
		LastParsed: time.Now(),
	}, nil
}

// NewService creates a new instance of the service with all dependencies.
func NewService(cr character.Repository) Service {
	return &service{
		characters: cr,
	}
}
