package statistics

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/nokka/d2-armory-api/internal/domain"
)

var validDifficulties = map[string]struct{}{
	domain.DifficultyNormal:    {},
	domain.DifficultyNightmare: {},
	domain.DifficultyHell:      {},
}

// Max data points is used to limit number of data points being returned
// since areas for example can host 138 entries.
const maxDataPoints = 8

// characterRepository is the interface representation of the data layer
// the service depend on.
type statisticsRepository interface {
	GetByCharacter(ctx context.Context, character string) (*domain.CharacterStatistics, error)
	Upsert(ctx context.Context, stat domain.StatisticsRequest) error
}

// Service performs all operations on statistics.
type Service struct {
	repository statisticsRepository
}

// GetCharacter will get the statistics on a specific character.
func (s Service) GetCharacter(ctx context.Context, character string) (*domain.CharacterStatistics, error) {
	char, err := s.repository.GetByCharacter(ctx, character)
	if err != nil {
		return nil, err
	}

	// Limit number of areas, to avoid showing all 138.
	if len(char.Normal.Area) > maxDataPoints {
		char.Normal.Area = getTopAreas((char.Normal.Area))
	}

	if len(char.Nightmare.Area) > maxDataPoints {
		char.Nightmare.Area = getTopAreas((char.Nightmare.Area))
	}

	if len(char.Hell.Area) > maxDataPoints {
		char.Hell.Area = getTopAreas((char.Hell.Area))
	}

	// Limit number of monsters, to avoid showing way too many.
	if len(char.Normal.Special) > maxDataPoints {
		char.Normal.Special = getTopSpecials(char.Normal.Special)
	}

	if len(char.Nightmare.Special) > maxDataPoints {
		char.Nightmare.Special = getTopSpecials(char.Nightmare.Special)
	}

	if len(char.Hell.Special) > maxDataPoints {
		char.Hell.Special = getTopSpecials(char.Hell.Special)
	}

	return char, nil
}

func getTopSpecials(monsters map[string]int) map[string]int {
	// Since maps are unordered by nature we need to temporarily
	// keep the kills in a struct and return a new map.
	type tmp struct {
		monster string
		kills   int
	}

	data := make([]tmp, 0)

	for monster, kills := range monsters {
		data = append(data, tmp{
			monster: monster,
			kills:   kills,
		})
	}

	sort.SliceStable(data, func(i, j int) bool {
		return data[i].kills > data[j].kills
	})

	topMonsters := map[string]int{}

	for _, v := range data[:maxDataPoints] {
		topMonsters[v.monster] = v.kills
	}

	return topMonsters
}

func getTopAreas(areas map[string]domain.AreaStats) map[string]domain.AreaStats {
	// Since maps are unordered by nature we need to temporarily
	// keep the areas in a struct and return a new map.
	type tmp struct {
		area  string
		kills uint
		time  uint
	}

	data := make([]tmp, 0)

	for area, vals := range areas {
		data = append(data, tmp{
			area:  area,
			kills: vals.Kills,
			time:  vals.Time,
		})
	}

	sort.SliceStable(data, func(i, j int) bool {
		return data[i].time > data[j].time
	})

	topAreas := map[string]domain.AreaStats{}

	for _, v := range data[:maxDataPoints] {
		topAreas[v.area] = domain.AreaStats{
			Time:  v.time,
			Kills: v.kills,
		}
	}

	return topAreas
}

// Parse will perform the storage of a statistics request.
func (s Service) Parse(ctx context.Context, stats []domain.StatisticsRequest) error {
	for _, req := range stats {
		if _, valid := validDifficulties[req.Difficulty]; !valid {
			return fmt.Errorf("difficulty %s supplied for character %s, %w", req.Difficulty, req.Character, domain.ErrRequest)
		}

		// Lower case the character name to keep consistency.
		req.Account = strings.ToLower(req.Account)
		req.Character = strings.ToLower(req.Character)

		// Upsert each character stat request.
		err := s.repository.Upsert(ctx, req)
		if err != nil {
			return err
		}
	}

	return nil
}

// NewService constructs a new statistics service with all the dependencies.
func NewService(repository statisticsRepository) *Service {
	return &Service{
		repository: repository,
	}
}
