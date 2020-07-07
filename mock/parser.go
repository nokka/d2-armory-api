package mock

import (
	"github.com/nokka/d2-armory-api/internal/domain"
)

// Parser is a mock implementation of the character parser.
type Parser struct {
	ParseFn      func(name string) (*domain.Character, error)
	ParseInvoked bool
}

// Parse parses the given binary name.
func (p Parser) Parse(name string) (*domain.Character, error) {
	p.ParseInvoked = true
	return p.ParseFn(name)
}
