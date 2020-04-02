package character

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/nokka/d2-armory-api/internal/domain"

	"github.com/nokka/d2-armory-api/mock"
)

func TestParseCharacter(t *testing.T) {
	type args struct {
		name          string
		ctx           context.Context
		cacheDuration time.Duration
	}

	type fields struct {
		characterRepository mock.CharacterRepository
	}

	tests := []struct {
		name   string
		args   args
		fields fields
	}{
		{
			name: "parse successful",
			args: args{
				name:          "nokka",
				ctx:           context.TODO(),
				cacheDuration: 1 * time.Minute,
			},
			fields: fields{
				characterRepository: mock.CharacterRepository{
					FindFn: func(ctx context.Context, id string) (*domain.Character, error) {
						return &domain.Character{}, nil
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService("/tmp", tt.fields.characterRepository, tt.args.cacheDuration)

			char, err := s.Parse(tt.args.ctx, tt.args.name)

			// TODO: Refactor the parser into it's own entity so it's mockable.
			fmt.Println(char, err)
		})
	}
}
