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
func (object *Object) Insert(eng *xorm.Engine) (affected int64, err error) {
	affected, err = eng.Insert(object)
	return
}

// Get -
func (object *Object) Get(eng *xorm.Engine) (has bool, err error) {
	has, err = eng.Get(object)
	return
}

// Update -
func (object *Object) Update(eng *xorm.Engine) (affected int64, err error) {
	affected, err = eng.ID(object.ID).Update(object)
	return
}

// AddTrait -
func (object *Object) AddTrait(eng *xorm.Engine, trait *Trait) (affected int64, err error) {
	affected, _, err = traitObjectInsertOrUpdate(eng, trait, object)
	return
}
