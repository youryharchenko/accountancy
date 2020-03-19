package accountancy

import (
	"time"

	//"github.com/go-xorm/xorm"
	"github.com/google/uuid"
	"xorm.io/xorm"
)

// Trait -
type Trait struct {
	ID      int64
	UUID    string                 `xorm:"varchar(36) notnull unique"`
	Name    string                 `xorm:"varchar(1024) notnull unique"`
	Props   map[string]interface{} `xorm:"jsonb"`
	Lib     string                 `xorm:"text"`
	Status  int                    `xorm:"notnull index"`
	Created time.Time              `xorm:"created"`
	Updated time.Time              `xorm:"updated"`
}

// TraitObject -
type TraitObject struct {
	ID       int64
	TraitID  int64                  `xorm:"notnull unique(trait_object)"`
	ObjectID int64                  `xorm:"notnull unique(trait_object)"`
	Props    map[string]interface{} `xorm:"jsonb"`
	Status   int                    `xorm:"notnull index"`
	Created  time.Time              `xorm:"created"`
	Updated  time.Time              `xorm:"updated"`
}

// NewTrait -
func NewTrait(name string, props map[string]interface{}, lib string) (trait *Trait, err error) {
	trait = &Trait{
		UUID:  uuid.New().String(),
		Name:  name,
		Props: props,
		Lib:   lib,
	}
	return
}

// Insert -
func (trait *Trait) Insert(eng *xorm.Engine) (affected int64, err error) {
	affected, err = eng.Insert(trait)
	return
}

// Get -
func (trait *Trait) Get(eng *xorm.Engine) (has bool, err error) {
	has, err = eng.Get(trait)
	return
}

// Update -
func (trait *Trait) Update(eng *xorm.Engine) (affected int64, err error) {
	affected, err = eng.ID(trait.ID).Update(trait)
	return
}

// AddRelationFrom -
func (trait *Trait) AddRelationFrom(eng *xorm.Engine, relation *Relation, traitFrom *Trait, props map[string]interface{}, lib string) (affected int64, err error) {
	affected, err = relationTraitsInsert(eng, relation, traitFrom, trait, props, lib)
	return
}

// AddRelationTo -
func (trait *Trait) AddRelationTo(eng *xorm.Engine, relation *Relation, traitTo *Trait, props map[string]interface{}, lib string) (affected int64, err error) {
	affected, err = relationTraitsInsert(eng, relation, trait, traitTo, props, lib)
	return
}

// AddObject -
func (trait *Trait) AddObject(eng *xorm.Engine, object *Object) (affected int64, err error) {
	affected, err = traitObjectInsert(eng, trait, object)
	return
}

// UpdateObjectProps -
func (trait *Trait) UpdateObjectProps(eng *xorm.Engine, object *Object, props map[string]interface{}) (affected int64, err error) {
	affected, err = traitObjectUpdate(eng, trait, object, props)
	return
}
