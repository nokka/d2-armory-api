package character

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/nokka/d2-armory-api/internal/domain"
)

func TestParseCharacter(t *testing.T) {
	type args struct {
		name          string
		ctx           context.Context
		cacheDuration time.Duration
	}

	type fields struct {
		characterRepository *characterRepositoryMock
		parser              *parserMock
	}

	type calls struct {
		storeCalls  int
		updateCalls int
		parseCalls  int
	}

	tests := []struct {
		name          string
		args          args
		fields        fields
		calls         calls
		expectedError error
	}{
		{
			name: "store successful",
			args: args{
				name:          "nokka",
				ctx:           context.TODO(),
				cacheDuration: 1 * time.Minute,
			},
			fields: fields{
				characterRepository: &characterRepositoryMock{
					FindFunc: func(ctx context.Context, id string) (*domain.Character, error) {
						return nil, domain.ErrNotFound
					},
					StoreFunc: func(ctx context.Context, character *domain.Character) error {
						return nil
					},
				},
				parser: &parserMock{
					ParseFunc: func(name string) (*domain.Character, error) {
						return &domain.Character{}, nil
					},
				},
			},
			calls: calls{
				storeCalls:  1,
				parseCalls:  1,
				updateCalls: 0,
			},
		},
		{
			name: "update successful",
			args: args{
				name:          "nokka",
				ctx:           context.TODO(),
				cacheDuration: 1 * time.Minute,
			},
			fields: fields{
				characterRepository: &characterRepositoryMock{
					FindFunc: func(ctx context.Context, id string) (*domain.Character, error) {
						return &domain.Character{}, nil
					},
					UpdateFunc: func(ctx context.Context, character *domain.Character) error {
						return nil
					},
				},
				parser: &parserMock{
					ParseFunc: func(name string) (*domain.Character, error) {
						return &domain.Character{}, nil
					},
				},
			},
			calls: calls{
				storeCalls:  0,
				parseCalls:  1,
				updateCalls: 1,
			},
		},
		{
			name: "temporary update error",
			args: args{
				name:          "nokka",
				ctx:           context.TODO(),
				cacheDuration: 1 * time.Minute,
			},
			fields: fields{
				characterRepository: &characterRepositoryMock{
					FindFunc: func(ctx context.Context, id string) (*domain.Character, error) {
						return &domain.Character{}, nil
					},
					UpdateFunc: func(ctx context.Context, character *domain.Character) error {
						return fmt.Errorf("temporary error: %w", domain.ErrTemporary)
					},
				},
				parser: &parserMock{
					ParseFunc: func(name string) (*domain.Character, error) {
						return &domain.Character{}, nil
					},
				},
			},
			calls: calls{
				storeCalls:  0,
				parseCalls:  1,
				updateCalls: 1,
			},
			expectedError: domain.ErrTemporary,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(tt.fields.parser, tt.fields.characterRepository, tt.args.cacheDuration)

			_, err := s.Parse(tt.args.ctx, tt.args.name)

			if err != nil && tt.expectedError == nil {
				t.Errorf("didn't expect an error, got = %v", err)
			}

			if tt.expectedError != nil && errors.Unwrap(err) != tt.expectedError {
				t.Errorf("Expected error to be = %v, got = %#v", tt.expectedError, errors.Unwrap(err))
			}

			if len(tt.fields.characterRepository.StoreCalls()) != tt.calls.storeCalls {
				t.Errorf("expected characterRepository.Store() to be called exactly %d times but was called %d times",
					tt.calls.storeCalls,
					len(tt.fields.characterRepository.StoreCalls()),
				)
			}

			if len(tt.fields.parser.ParseCalls()) != tt.calls.parseCalls {
				t.Errorf("expected parser.Parse() to be called exactly %d times but was called %d times",
					tt.calls.parseCalls,
					len(tt.fields.parser.ParseCalls()),
				)
			}

			if len(tt.fields.characterRepository.UpdateCalls()) != tt.calls.updateCalls {
				t.Errorf("expected characterRepository.UpdateCalls() to be called exactly %d times but was called %d times",
					tt.calls.updateCalls,
					len(tt.fields.characterRepository.UpdateCalls()),
				)
			}
		})
	}
}
