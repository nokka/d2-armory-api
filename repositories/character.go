package repositories

import (
	"github.com/nokka/armory/character"
	"gopkg.in/mgo.v2"
)

type characterRepository struct {
	db      string
	session *mgo.Session
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
