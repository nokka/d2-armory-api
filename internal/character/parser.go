package character

import (
	"fmt"
	"os"
	"time"

	"github.com/nokka/d2-armory-api/internal/domain"
	"github.com/nokka/d2s"
)

func parseCharacter(name, path string) (*domain.Character, error) {
	// Character path on disk.
	file, err := os.Open(fmt.Sprintf("%s/%s", path, name))
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
