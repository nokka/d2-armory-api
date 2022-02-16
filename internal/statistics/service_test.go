package statistics

import (
	"context"
	"errors"
	"testing"

	"github.com/nokka/d2-armory-api/internal/domain"
)

func TestParse(t *testing.T) {
	type args struct {
		ctx   context.Context
		stats []domain.StatisticsRequest
	}

	type fields struct {
		statisticsRepository *statisticsRepositoryMock
	}

	type calls struct {
		upsertCalls int
	}

	tests := []struct {
		name          string
		args          args
		fields        fields
		calls         calls
		expectedError bool
	}{
		{
			name: "parse successful",
			args: args{
				ctx: context.TODO(),
				stats: []domain.StatisticsRequest{
					{Character: "nokka", Difficulty: domain.DifficultyHell},
				},
			},
			fields: fields{
				statisticsRepository: &statisticsRepositoryMock{
					UpsertFunc: func(ctx context.Context, stat domain.StatisticsRequest) error {
						return nil
					},
				},
			},
			calls: calls{
				upsertCalls: 1,
			},
		},
		{
			name: "parse error",
			args: args{
				ctx: context.TODO(),
				stats: []domain.StatisticsRequest{
					{Character: "nokka", Difficulty: domain.DifficultyNormal},
				},
			},
			fields: fields{
				statisticsRepository: &statisticsRepositoryMock{
					UpsertFunc: func(ctx context.Context, stat domain.StatisticsRequest) error {
						return errors.New("something went wrong")
					},
				},
			},
			calls: calls{
				upsertCalls: 1,
			},
			expectedError: true,
		},
		{
			name: "invalid difficulty supplied",
			args: args{
				ctx: context.TODO(),
				stats: []domain.StatisticsRequest{
					{Character: "nokka", Difficulty: "invalid"},
				},
			},
			fields: fields{
				statisticsRepository: &statisticsRepositoryMock{},
			},
			calls: calls{
				upsertCalls: 0,
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(tt.fields.statisticsRepository)

			err := s.Parse(tt.args.ctx, tt.args.stats)

			if err != nil && tt.expectedError == false {
				t.Errorf("didn't expect an error, got = %v", err)
			}

			if len(tt.fields.statisticsRepository.UpsertCalls()) != tt.calls.upsertCalls {
				t.Errorf("expected statisticsRepository.UpsertCalls() to be called exactly %d times but was called %d times",
					tt.calls.upsertCalls,
					len(tt.fields.statisticsRepository.UpsertCalls()),
				)
			}
		})
	}
}

func TestGetCharacter(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
	}

	type fields struct {
		statisticsRepository *statisticsRepositoryMock
	}

	tests := []struct {
		name             string
		args             args
		fields           fields
		expectedError    bool
		normalAreas      int
		excludedAreas    []string
		normalSpecial    int
		excludedSpecials []string
	}{
		{
			name: "get character successful",
			args: args{
				ctx:  context.TODO(),
				name: "nokka",
			},
			fields: fields{
				statisticsRepository: &statisticsRepositoryMock{
					GetByCharacterFunc: func(ctx context.Context, character string) (*domain.CharacterStatistics, error) {
						return &domain.CharacterStatistics{
							Normal: domain.Stats{
								Area: map[string]domain.AreaStats{
									"Stony field": {Kills: 1, Time: 120},
								},
							},
						}, nil
					},
				},
			},
			normalAreas: 1,
		},
		{
			name: "get character exclude areas",
			args: args{
				ctx:  context.TODO(),
				name: "nokka",
			},
			fields: fields{
				statisticsRepository: &statisticsRepositoryMock{
					GetByCharacterFunc: func(ctx context.Context, character string) (*domain.CharacterStatistics, error) {
						return &domain.CharacterStatistics{
							Normal: domain.Stats{
								Area: map[string]domain.AreaStats{
									"Burial grounds":     {Kills: 1, Time: 90},
									"Cold plains":        {Kills: 1, Time: 80},
									"Moo Moo farm":       {Kills: 1, Time: 70},
									"The Pit level 2":    {Kills: 1, Time: 60},
									"Blood moor":         {Kills: 1, Time: 5},
									"River of flame":     {Kills: 1, Time: 50},
									"The sewers":         {Kills: 1, Time: 40},
									"Stony field":        {Kills: 1, Time: 10},
									"Worldstone chamber": {Kills: 1, Time: 30},
									"Tal rashas tomb":    {Kills: 1, Time: 20},
								},
							},
						}, nil
					},
				},
			},
			normalAreas: 8,
			excludedAreas: []string{
				"Stony field",
				"Blood moor",
			},
		},
		{
			name: "get character exclude special monsters",
			args: args{
				ctx:  context.TODO(),
				name: "nokka",
			},
			fields: fields{
				statisticsRepository: &statisticsRepositoryMock{
					GetByCharacterFunc: func(ctx context.Context, character string) (*domain.CharacterStatistics, error) {
						return &domain.CharacterStatistics{
							Normal: domain.Stats{
								Special: map[string]int{
									"Andariel":    42,
									"Baal":        13,
									"Corpsefire":  12,
									"Coldcrow":    11,
									"Diablo":      30,
									"Duriel":      100,
									"Mephisto":    20,
									"Nihlatalak":  202,
									"Blood Raven": 202,
									"Griswold":    1,
								},
							},
						}, nil
					},
				},
			},
			normalSpecial: 8,
			excludedSpecials: []string{
				"Griswold",
				"Coldcrow",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(tt.fields.statisticsRepository)

			stats, err := s.GetCharacter(tt.args.ctx, tt.args.name)

			if err != nil && tt.expectedError == false {
				t.Errorf("didn't expect an error, got = %v", err)
			}

			if len(stats.Normal.Area) != tt.normalAreas {
				t.Errorf("expected number of areas to be = %d, got = %d", tt.normalAreas, len(stats.Normal.Area))
			}

			for _, name := range tt.excludedAreas {
				if _, exists := stats.Normal.Area[name]; exists {
					t.Errorf("did not expect area %s to still be in the map", name)
				}
			}

			for _, name := range tt.excludedSpecials {
				if _, exists := stats.Normal.Special[name]; exists {
					t.Errorf("did not expect special monster %s to still be in the map", name)
				}
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type args struct {
		ctx       context.Context
		character string
	}

	type fields struct {
		statisticsRepository *statisticsRepositoryMock
	}

	tests := []struct {
		name          string
		args          args
		fields        fields
		deleteCalls   int
		expectedError bool
	}{
		{
			name: "delete character stats successfully",
			args: args{
				ctx:       context.TODO(),
				character: "nokka",
			},
			fields: fields{
				statisticsRepository: &statisticsRepositoryMock{
					DeleteFunc: func(ctx context.Context, character string) error {
						return nil
					},
				},
			},
			deleteCalls:   1,
			expectedError: false,
		},
		{
			name: "delete character name too short",
			args: args{
				ctx:       context.TODO(),
				character: "n",
			},
			fields: fields{
				statisticsRepository: &statisticsRepositoryMock{},
			},
			deleteCalls:   0,
			expectedError: true,
		},
		{
			name: "delete character returns error",
			args: args{
				ctx:       context.TODO(),
				character: "nokka",
			},
			fields: fields{
				statisticsRepository: &statisticsRepositoryMock{
					DeleteFunc: func(ctx context.Context, character string) error {
						return errors.New("something went terribly wrong")
					},
				},
			},
			deleteCalls:   1,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(tt.fields.statisticsRepository)

			err := s.DeleteStats(tt.args.ctx, tt.args.character)

			if (err != nil) != tt.expectedError {
				t.Errorf("got error = %v, expectedError %v", err, tt.expectedError)
			}

			deleteCalls := len(tt.fields.statisticsRepository.DeleteCalls())
			if deleteCalls != tt.deleteCalls {
				t.Errorf("Delete() was called %d times, expected to be called %d times", deleteCalls, tt.deleteCalls)
			}
		})
	}
}
