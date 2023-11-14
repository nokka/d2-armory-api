package mgo

import (
	"context"
	"fmt"
	"strings"

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
func (r *StatisticsRepository) GetByCharacter(ctx context.Context, character string) (*domain.CharacterStatistics, error) {
	var char domain.CharacterStatistics
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

	difficulty := strings.ToLower(stat.Difficulty)

	// Difficulty updates.
	values := map[string]interface{}{
		fmt.Sprintf("%s.total_kills", difficulty):        stat.TotalKills,
		fmt.Sprintf("%s.total_unique_kills", difficulty): stat.TotalUniqueKills,
		fmt.Sprintf("%s.total_champ_kills", difficulty):  stat.TotalChampKills,
	}

	// Looping over special monsters to add them to upsert.
	for monster, val := range stat.Special {
		values[fmt.Sprintf("%s.special.%s", difficulty, monster)] = val
	}

	// Looping over area statistics to add them to upsert.
	for area, val := range stat.Area {
		values[fmt.Sprintf("%s.area.%s.kills", difficulty, area)] = val.Kills
		values[fmt.Sprintf("%s.area.%s.time", difficulty, area)] = val.Time
		values[fmt.Sprintf("%s.area.%s.unique_kills", difficulty, area)] = val.UniqueKills
		values[fmt.Sprintf("%s.area.%s.champ_kills", difficulty, area)] = val.ChampKills
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

// Delete will delete statistics about the given character.
func (r *StatisticsRepository) Delete(ctx context.Context, character string) error {
	// Query to find document on.
	query := bson.M{
		"character": character,
	}

	_, err := r.client.Database(r.db).Collection(statCollectionName).
		DeleteOne(ctx, query, nil)
	if err != nil {
		return mongoErr(err)
	}

	return nil
}

// Internal store function to create the document for the first time.
func (r *StatisticsRepository) store(ctx context.Context, request domain.StatisticsRequest) error {
	// Initiate character document, be explicit about the maps to avoid upsert errors on nil.
	character := domain.CharacterStatistics{
		Account:   request.Account,
		Character: request.Character,
		Normal: domain.Stats{
			Special: make(map[string]int, 0),
			Area:    make(map[string]domain.AreaStats, 0),
		},
		Nightmare: domain.Stats{
			Special: make(map[string]int, 0),
			Area:    make(map[string]domain.AreaStats, 0),
		},
		Hell: domain.Stats{
			Special: make(map[string]int, 0),
			Area:    make(map[string]domain.AreaStats, 0),
		},
	}

	stats := domain.Stats{
		TotalKills:       request.TotalKills,
		TotalUniqueKills: request.TotalUniqueKills,
		TotalChampKills:  request.TotalChampKills,
		Special:          request.Special,
		Area:             request.Area,
	}

	switch request.Difficulty {
	case domain.DifficultyNormal:
		character.Normal = stats
	case domain.DifficultyNightmare:
		character.Nightmare = stats
	case domain.DifficultyHell:
		character.Hell = stats
	}

	_, err := r.client.Database(r.db).Collection(statCollectionName).
		InsertOne(ctx, character)
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
