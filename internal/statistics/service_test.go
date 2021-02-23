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
					{Account: "nokka"},
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
					{Account: "nokka"},
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
