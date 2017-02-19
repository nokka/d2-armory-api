// Package retrieving is responsible for retrieving and parsing character data
// from the diablo character d2s files.
package retrieving

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/nokka/armory/character"
	"github.com/nokka/d2s"
)

// ErrInvalidArgument is returned when one or more arguments are invalid.
var ErrInvalidArgument = errors.New("invalid argument")

// ErrNonExistingCharacter is returned when someone requests data for a d2s char
// that does not exists on disk.
var ErrNonExistingCharacter = errors.New("The requested character does not exist")

// ErrFailedToParse is returned when someone requests a character and the parser
// return an error.
var ErrFailedToParse = errors.New("The requested character could not be parsed")

// The name regexp required for character names, to enforce strict diablo rules
// on the names to prevent missuse of the endpoint.
const nameRegexp = "^[a-zA-Z]+[_-]?[a-zA-Z]+$"

// Service provides operations on d2s character data.
type Service interface {
	// GetCharacter will return the character with the given name.
	RetrieveCharacter(string) (*character.Character, error)
}

type service struct {
	characters character.Repository
	d2spath    string
}

func (s *service) RetrieveCharacter(name string) (*character.Character, error) {
	if name == "" {
		return nil, ErrInvalidArgument
	}

	match, _ := regexp.MatchString(nameRegexp, name)
	if !match {
		return nil, ErrInvalidArgument
	}

	// Find character in collection.
	c := s.characters.Find(name)

	if c != nil {
		// Check the time when we last parsed it.
		diff := time.Since(c.LastParsed)

		// If we haven't parsed this char in the last 10 minutes, lets parse it.
		if diff.Minutes() >= 10 {
			parsed, err := s.parseCharacter(name)
			if err != nil {
				return nil, err
			}

			// Update the existing record in the db.
			err = s.characters.Update(parsed)
			if err != nil {
				return nil, err
			}

			return parsed, nil
		}

		// We parsed this character recently, lets return the copy from the db.
		return c, nil
	}

	// Character didn't exist at all, so lets parse and store it.
	parsed, err := s.parseCharacter(name)
	if err != nil {
		return nil, err
	}

	if err := s.characters.Store(parsed); err != nil {
		return nil, err
	}

	return parsed, nil
}

func (s *service) parseCharacter(name string) (*character.Character, error) {
	path := s.d2spath + name
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return nil, ErrNonExistingCharacter
	}

	char, err := d2s.Parse(file)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	c := character.Character{
		ID:         name,
		D2s:        &char,
		LastParsed: time.Now(),
	}

	return &c, nil
}

// NewService creates a new instance of the service with all dependencies.
func NewService(cr character.Repository, d2spath string) Service {
	return &service{
		characters: cr,
		d2spath:    d2spath,
	}
}
