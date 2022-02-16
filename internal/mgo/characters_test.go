//go:build integration

package mgo

import (
	"context"
	"testing"
	"time"

	"github.com/nokka/d2-armory-api/internal/domain"
	"github.com/nokka/d2-armory-api/pkg/env"
	"github.com/nokka/d2s"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func TestCharacterRepository(t *testing.T) {
	clientOptions := options.Client().ApplyURI("mongodb://" + env.String("MONGO_HOST", "mongodb:27017"))

	clientOptions.SetAuth(options.Credential{
		AuthSource: env.String("MONGO_DB", "armory"),
		Username:   env.String("MONGO_USERNAME", "armory"),
		Password:   env.String("MONGO_PASSWORD", "not_secure_at_all"),
	})

	// Context used for mongo operations, to time them out and cancel their context.
	mgoCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(mgoCtx, clientOptions)
	if err != nil {
		t.Error("failed to connect to mongodb", err)
		return
	}

	err = client.Ping(mgoCtx, readpref.Primary())
	if err != nil {
		t.Error("failed to ping mongodb", err)
		return
	}

	characterRepository := NewCharacterRepository("armory", client)

	t.Run("store character", func(t *testing.T) {
		err := characterRepository.Store(mgoCtx, &domain.Character{
			ID:         "nokka",
			D2s:        &d2s.Character{},
			LastParsed: time.Now(),
		})
		if err != nil {
			t.Error("failed to store character")
		}
	})

	t.Run("update character", func(t *testing.T) {
		err := characterRepository.Update(mgoCtx, &domain.Character{
			ID:         "nokka",
			D2s:        &d2s.Character{},
			LastParsed: time.Now(),
		})
		if err != nil {
			t.Error("failed to update character")
		}
	})

	t.Run("find character by id", func(t *testing.T) {
		character, err := characterRepository.Find(mgoCtx, "nokka")
		if err != nil {
			t.Error("failed to get character")
		}

		if character.ID != "nokka" {
			t.Error("failed to get character by the ID")
		}
	})
}
