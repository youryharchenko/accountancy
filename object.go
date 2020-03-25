package accountancy

import (
	"time"

	"github.com/google/uuid"
	"xorm.io/xorm"
)

// Object -
type Object struct {
	ID      int64
	UUID    string                 `xorm:"varchar(36) notnull unique"`
	Name    string                 `xorm:"varchar(1024) notnull unique"`
	Props   map[string]interface{} `xorm:"jsonb"`
	Status  int                    `xorm:"notnull index"`
	Created time.Time              `xorm:"created"`
	Updated time.Time              `xorm:"updated"`
}

// NewObject -
func NewObject(name string, props map[string]interface{}) (object *Object, err error) {
	object = &Object{
		UUID:  uuid.New().String(),
		Name:  name,
		Props: props,
	}
	return
}

// Insert -
func (object *Object) Insert(db DB) (affected int64, err error) {
	affected, err = db.Insert(object)
	return
}

// Get -
func (object *Object) Get(db DB) (has bool, err error) {
	has, err = db.Get(object)
	return
}

// Update -
func (object *Object) Update(db DB) (affected int64, err error) {
	affected, err = db.ID(object.ID).Update(object)
	return
}

// InsertOrUpdate -
func (object *Object) InsertOrUpdate(db DB, find *Object) (affected int64, inserted bool, err error) {

	ok, err := find.Get(db)
	if err != nil {
		return
	}

	if ok {
		object.ID = find.ID
		_, err = object.Update(db)
		if err != nil {
			return
		}
		inserted = false
	} else {
		_, err = object.Insert(db)
		if err != nil {
			return
		}
		inserted = true
	}

	return
}

// AddTrait -
func (object *Object) AddTrait(eng *xorm.Engine, trait *Trait) (affected int64, err error) {
	affected, _, err = traitObjectInsertOrUpdate(eng, trait, object)
	return
}
