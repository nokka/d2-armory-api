package character

import (
	"fmt"
	"os"

	"github.com/nokka/d2-armory-api/internal/domain"
	"github.com/nokka/d2s"
)

func parseCharacter(path string) (*d2s.Character, error) {
	fmt.Println(path)
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("character binary does not exist: %w", domain.ErrNotFound)
	}

	// Close the file when we're done.
	defer file.Close()

	// Parse the actual .d2s binary file.
	char, err := d2s.Parse(file)
	if err != nil {
		return nil, fmt.Errorf("binary parse error: %w", err)
	}

	return char, nil
}
