package character

import (
	"context"
	"errors"
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
		parser              mock.Parser
	}

	tests := []struct {
		name          string
		args          args
		fields        fields
		expectedError error
		storeInvoked  bool
		updateInvoked bool
	}{
		{
			name: "store successful",
			args: args{
				name:          "nokka",
				ctx:           context.TODO(),
				cacheDuration: 1 * time.Minute,
			},
			fields: fields{
				characterRepository: mock.CharacterRepository{
					FindFn: func(ctx context.Context, id string) (*domain.Character, error) {
						return nil, domain.ErrNotFound
					},
					StoreFn: func(ctx context.Context, character *domain.Character) error {
						return nil
					},
				},
				parser: mock.Parser{
					ParseFn: func(name string) (*domain.Character, error) {
						return &domain.Character{}, nil
					},
				},
			},
			storeInvoked: true,
		},
		{
			name: "update successful",
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
					UpdateFn: func(ctx context.Context, character *domain.Character) error {
						return nil
					},
				},
				parser: mock.Parser{
					ParseFn: func(name string) (*domain.Character, error) {
						return &domain.Character{}, nil
					},
				},
			},
			storeInvoked:  false,
			updateInvoked: true,
		},
		{
			name: "temporary update error",
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
					UpdateFn: func(ctx context.Context, character *domain.Character) error {
						return fmt.Errorf("temporary error: %w", domain.ErrTemporary)
					},
				},
				parser: mock.Parser{
					ParseFn: func(name string) (*domain.Character, error) {
						return &domain.Character{}, nil
					},
				},
			},
			storeInvoked:  false,
			updateInvoked: true,
			expectedError: domain.ErrTemporary,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(tt.fields.parser, &tt.fields.characterRepository, tt.args.cacheDuration)

			_, err := s.Parse(tt.args.ctx, tt.args.name)

			if err != nil && tt.expectedError == nil {
				t.Errorf("didn't expect an error, got = %v", err)
			}

			if tt.expectedError != nil && errors.Unwrap(err) != tt.expectedError {
				t.Errorf("Expected error to be = %v, got = %#v", tt.expectedError, errors.Unwrap(err))
			}

			if tt.fields.characterRepository.StoreInvoked != tt.storeInvoked {
				t.Errorf("expected Store() invocation to be %v", tt.storeInvoked)
			}

			if tt.fields.characterRepository.UpdateInvoked != tt.updateInvoked {
				t.Errorf("expected Update() invocation to be %v", tt.updateInvoked)
			}
		})
	}
}
