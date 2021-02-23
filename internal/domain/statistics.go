package domain

// StatisticsRequest represents a request to store statistics.
type StatisticsRequest struct {
	Account    string         `json:"account"`
	Character  string         `json:"character"`
	Uniques    int            `json:"uniques"`
	Champions  int            `json:"champions"`
	TotalKills int            `json:"totalkills"`
	Special    map[string]int `json:"special"`
	Regular    map[string]int `json:"regular"`
}
