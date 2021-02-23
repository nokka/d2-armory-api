package mgo

import (
	"context"
	"fmt"

	"github.com/nokka/d2-armory-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	// statCollectionName is the name of the collection we'll use for all queries.
	statCollectionName = "statistics"
)

// StatisticsRepository handles all operations on statistics.
type StatisticsRepository struct {
	db     string
	client *mongo.Client
}

// GetByCharacter will return statistics for the character.
func (r *StatisticsRepository) GetByCharacter(ctx context.Context, character string) (*domain.StatisticsRequest, error) {
	var char domain.StatisticsRequest

	err := r.client.Database(r.db).Collection(statCollectionName).
		FindOne(ctx, bson.M{"character": character}).Decode(&char)
	if err != nil {
		return nil, mongoErr(err)
	}

	return &char, nil
}

// Upsert will upsert statistics about the given character.
func (r *StatisticsRepository) Upsert(ctx context.Context, stat domain.StatisticsRequest) error {
	// Query to find document on.
	query := bson.M{
		"character": stat.Character,
	}

	// Root document updates.
	values := map[string]interface{}{
		"champions":  stat.Champions,
		"uniques":    stat.Uniques,
		"totalkills": stat.TotalKills,
	}

	// Looping over special monsters to add them to upsert.
	for monster, val := range stat.Special {
		values[fmt.Sprintf("special.%s", monster)] = val
	}

	// Looping over regular monsters to add them to upsert.
	for monster, val := range stat.Regular {
		values[fmt.Sprintf("regular.%s", monster)] = val
	}

	result, err := r.client.Database(r.db).Collection(statCollectionName).
		UpdateOne(ctx, query, bson.M{"$inc": values})
	if err != nil {
		return mongoErr(err)
	}

	// First time the character is indexed, just store it.
	if result.MatchedCount == 0 {
		return r.store(ctx, stat)
	}

	return nil
}

// Internal store function to create the document for the first time.
func (r *StatisticsRepository) store(ctx context.Context, stat domain.StatisticsRequest) error {
	_, err := r.client.Database(r.db).Collection(statCollectionName).
		InsertOne(ctx, stat)
	if err != nil {
		return mongoErr(err)
	}

	return nil
}

// NewStatisticsRepository returns a new instance of a MongoDB statistics repository.
func NewStatisticsRepository(db string, client *mongo.Client) *StatisticsRepository {
	return &StatisticsRepository{
		db:     db,
		client: client,
	}
}
