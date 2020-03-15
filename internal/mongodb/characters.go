package mongodb

import (
	"fmt"
	"time"

	"github.com/nokka/d2-armory-api/internal/domain"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// CharacterRepository handles all operations on characters.
type CharacterRepository struct {
	db        string
	connector *Connector
}

// Find will find a character by name.
func (r *CharacterRepository) Find(id string) (*domain.Character, error) {
	session, err := r.getSession()
	if err != nil {
		return nil, err
	}

	c := session.DB(r.db).C("character")

	var result domain.Character
	if err := c.Find(bson.M{"id": id}).One(&result); err != nil {
		return nil, mongoErr(err)
	}

	fmt.Println("RETURNING RESULT")

	return &result, nil
}

// Update ...
func (r *CharacterRepository) Update(character *domain.Character) error {
	session, err := r.getSession()
	if err != nil {
		return err
	}

	c := session.DB(r.db).C("character")
	change := bson.M{"$set": bson.M{"d2s": character.D2s, "lastparsed": time.Now()}}

	return c.Update(bson.M{"id": character.ID}, change)
}

// Store ...
func (r *CharacterRepository) Store(character *domain.Character) error {
	session, err := r.getSession()
	if err != nil {
		return err
	}

	c := session.DB(r.db).C("character")
	return c.Insert(character)
}

// getSession will return a new session if it has not been set yet.
func (r *CharacterRepository) getSession() (*mgo.Session, error) {
	s := r.connector.GetSession()
	if s == nil {
		return nil, fmt.Errorf("unable to connect to database: %w", domain.ErrUnavailable)
	}

	return s, nil
}

// NewCharacterRepository returns a new instance of a MongoDB character repository.
func NewCharacterRepository(db string, connector *Connector) *CharacterRepository {
	return &CharacterRepository{
		db:        db,
		connector: connector,
	}
}
