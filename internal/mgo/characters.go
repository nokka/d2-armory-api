package mgo

import (
	"context"
	"time"

	"github.com/nokka/d2-armory-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	// characterCollectionName is the name of the collection we'll use for all queries.
	characterCollectionName = "character"
)

// CharacterRepository handles all operations on characters.
type CharacterRepository struct {
	db     string
	client *mongo.Client
}

// Find will find a character by name.
func (r *CharacterRepository) Find(ctx context.Context, id string) (*domain.Character, error) {
	// Struct to decode query result into.
	var char domain.Character

	// Find the character by id in the collection.
	err := r.client.Database(r.db).Collection(characterCollectionName).
		FindOne(ctx, bson.M{"id": id}).Decode(&char)
	if err != nil {
		return nil, mongoErr(err)
	}

	return &char, nil
}

// Update will update the given resource.
func (r *CharacterRepository) Update(ctx context.Context, character *domain.Character) error {
	// Changeset, update the binary and time of parsing.
	change := bson.M{
		"$set": bson.M{
			"d2s":        character.D2s,
			"lastparsed": time.Now(),
		},
	}

	_, err := r.client.Database(r.db).Collection(characterCollectionName).
		UpdateOne(ctx, bson.M{"id": character.ID}, change)
	if err != nil {
		return mongoErr(err)
	}

	return nil
}

// Store will store the new resource.
func (r *CharacterRepository) Store(ctx context.Context, character *domain.Character) error {
	_, err := r.client.Database(r.db).Collection(characterCollectionName).
		InsertOne(ctx, character)
	if err != nil {
		return mongoErr(err)
	}

	return nil
}

// NewCharacterRepository returns a new instance of a MongoDB character repository.
func NewCharacterRepository(db string, client *mongo.Client) *CharacterRepository {
	return &CharacterRepository{
		db:     db,
		client: client,
	}
}
