package parsing

import (
	"fmt"
	"os"
	"time"

	"github.com/nokka/d2-armory-api/internal/domain"
	"github.com/nokka/d2s"
)

// Parser performs all parsing from d2s data to our domain model.
type Parser struct {
	d2spath string
}

// Parse will parse the given character on disk into a character in our domain model.
func (p Parser) Parse(name string) (*domain.Character, error) {
	// Character path on disk.
	file, err := os.Open(fmt.Sprintf("%s/%s", p.d2spath, name))
	if err != nil {
		return nil, fmt.Errorf("character binary does not exist: %w", domain.ErrNotFound)
	}

	// Close the file when we're done.
	defer file.Close()

	// Parse the actual .d2s binary file.
	d2schar, err := d2s.Parse(file)
	if err != nil {
		return nil, fmt.Errorf("binary parse error: %w", err)
	}

	character := domain.Character{
		ID:         name,
		D2s:        d2schar,
		LastParsed: time.Now(),
	}

	return &character, nil
}

// NewParser constructs a new parser with dependencies.
func NewParser(d2spath string) *Parser {
	return &Parser{
		d2spath: d2spath,
	}
}
