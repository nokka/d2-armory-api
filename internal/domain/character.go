package domain

import (
	"time"

	"github.com/nokka/d2s"
)

// Character represents a Diablo II character.
type Character struct {
	ID         string         `json:"d2s_id"`
	D2s        *d2s.Character `json:"d2s"`
	LastParsed time.Time      `json:"last_parsed"`
}
