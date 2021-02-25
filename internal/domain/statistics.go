package domain

// Difficulty constants available in the statistics.
const (
	DifficultyNormal    = "Normal"
	DifficultyNightmare = "Nightmare"
	DifficultyHell      = "Hell"
)

// StatisticsRequest represents statistics request.
type StatisticsRequest struct {
	Account    string         `json:"account"`
	Character  string         `json:"character"`
	Difficulty string         `json:"difficulty"`
	Uniques    int            `json:"uniques"`
	Champions  int            `json:"champions"`
	TotalKills int            `json:"totalkills"`
	Special    map[string]int `json:"special"`
	Regular    map[string]int `json:"regular"`
}

// CharacterStatistics represents character statistics.
type CharacterStatistics struct {
	Account   string `json:"account"`
	Character string `json:"character"`
	Normal    Stats  `json:"normal"`
	Nightmare Stats  `json:"nightmare"`
	Hell      Stats  `json:"hell"`
}

// Stats is repeated for each difficulty.
type Stats struct {
	Uniques    int            `json:"uniques"`
	Champions  int            `json:"champions"`
	TotalKills int            `json:"total_kills" bson:"total_kills"`
	Special    map[string]int `json:"special"`
	Regular    map[string]int `json:"regular"`
}
