package repositories

import (
	"time"

	"github.com/nokka/armory/character"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type characterRepository struct {
	db      string
	session *mgo.Session
}

func (r *characterRepository) Find(id string) *character.Character {

	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.db).C("character")

	var result *character.Character
	if err := c.Find(bson.M{"id": id}).One(&result); err != nil {
		return nil
	}

	return result
}

func (r *characterRepository) Update(character *character.Character) error {
	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.db).C("character")
	change := bson.M{"$set": bson.M{"d2s": character.D2s, "lastparsed": time.Now()}}
	err := c.Update(bson.M{"id": character.ID}, change)
	if err != nil {
		panic(err)
	}

	return err
}

func (r *characterRepository) Store(character *character.Character) error {
	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.db).C("character")
	err := c.Insert(character)

	return err
}

// NewCharacterRepository returns a new instance of a MongoDB character repository.
func NewCharacterRepository(db string, session *mgo.Session) (character.Repository, error) {
	r := &characterRepository{
		db:      db,
		session: session,
	}

	return r, nil
}
