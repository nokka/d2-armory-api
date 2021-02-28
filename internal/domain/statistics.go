package domain

// Difficulty constants available in the statistics.
const (
	DifficultyNormal    = "Normal"
	DifficultyNightmare = "Nightmare"
	DifficultyHell      = "Hell"
)

// StatisticsRequest represents statistics request.
type StatisticsRequest struct {
	Account    string               `json:"account"`
	Character  string               `json:"character"`
	Difficulty string               `json:"difficulty"`
	TotalKills int                  `json:"totalkills"`
	Special    map[string]int       `json:"special"`
	Area       map[string]AreaStats `json:"area"`
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
	TotalKills int                  `json:"total_kills" bson:"total_kills"`
	Special    map[string]int       `json:"special"`
	Area       map[string]AreaStats `json:"area"`
}

// AreaStats contains information about a particular area.
type AreaStats struct {
	Kills uint `json:"kills"`
	Time  uint `json:"time"`
}
