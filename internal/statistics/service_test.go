package statistics

import (
	"context"
	"errors"
	"testing"

	"github.com/nokka/d2-armory-api/internal/domain"

	"github.com/nokka/d2-armory-api/mock"
)

func TestParse(t *testing.T) {
	type args struct {
		name  string
		ctx   context.Context
		stats []domain.StatisticsRequest
	}

	type fields struct {
		statisticsRepository mock.StatisticsRepository
	}

	tests := []struct {
		name          string
		args          args
		fields        fields
		expectedError bool
		upsertInvoked bool
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
				statisticsRepository: mock.StatisticsRepository{
					UpsertFn: func(ctx context.Context, stat domain.StatisticsRequest) error {
						return nil
					},
				},
			},
			upsertInvoked: true,
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
				statisticsRepository: mock.StatisticsRepository{
					UpsertFn: func(ctx context.Context, stat domain.StatisticsRequest) error {
						return errors.New("something went wrong")
					},
				},
			},
			upsertInvoked: true,
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
			fields:        fields{},
			upsertInvoked: false,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(&tt.fields.statisticsRepository)

			err := s.Parse(tt.args.ctx, tt.args.stats)

			if err != nil && tt.expectedError == false {
				t.Errorf("didn't expect an error, got = %v", err)
			}

			if tt.fields.statisticsRepository.UpsertInvoked != tt.upsertInvoked {
				t.Errorf("expected Upsert() invocation to be %v", tt.upsertInvoked)
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
		statisticsRepository mock.StatisticsRepository
	}

	tests := []struct {
		name          string
		args          args
		fields        fields
		expectedError bool
		normalAreas   int
		excludedAreas []string
	}{
		{
			name: "get character successful",
			args: args{
				ctx:  context.TODO(),
				name: "nokka",
			},
			fields: fields{
				statisticsRepository: mock.StatisticsRepository{
					GetByCharacterFn: func(ctx context.Context, character string) (*domain.CharacterStatistics, error) {
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
				statisticsRepository: mock.StatisticsRepository{
					GetByCharacterFn: func(ctx context.Context, character string) (*domain.CharacterStatistics, error) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(&tt.fields.statisticsRepository)

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
		})
	}
}
