package accountancy

import (
	"time"

	"github.com/go-xorm/xorm"
	"github.com/google/uuid"
)

// Relation -
type Relation struct {
	ID      int64
	UUID    string                 `xorm:"varchar(36) notnull unique"`
	Name    string                 `xorm:"varchar(1024) notnull unique"`
	Props   map[string]interface{} `xorm:"jsonb"`
	Status  int                    `xorm:"notnull index"`
	Created time.Time              `xorm:"created"`
	Updated time.Time              `xorm:"updated"`
}

// RelationTraits -
type RelationTraits struct {
	ID          int64
	RelationID  int64                  `xorm:"notnull unique(relation_traits)"`
	TraitFromID int64                  `xorm:"notnull unique(relation_traits)"`
	TraitToID   int64                  `xorm:"notnull unique(relation_traits)"`
	Props       map[string]interface{} `xorm:"jsonb"`
	Lib         string                 `xorm:"text"`
	Status      int                    `xorm:"notnull index"`
	Created     time.Time              `xorm:"created"`
	Updated     time.Time              `xorm:"updated"`
}

// NewRelation -
func NewRelation(name string, props map[string]interface{}, lib string) (relation *Relation, err error) {
	relation = &Relation{
		UUID:  uuid.New().String(),
		Name:  name,
		Props: props,
	}
	return
}

// Insert -
func (relation *Relation) Insert(eng *xorm.Engine) (affected int64, err error) {
	affected, err = eng.Insert(relation)
	return
}

// Get -
func (relation *Relation) Get(eng *xorm.Engine) (has bool, err error) {
	has, err = eng.Get(relation)
	return
}

// Update -
func (relation *Relation) Update(eng *xorm.Engine) (affected int64, err error) {
	affected, err = eng.ID(relation.ID).Update(relation)
	return
}

// AddRelationTraits -
func (relation *Relation) AddRelationTraits(eng *xorm.Engine, traitFrom *Trait, traitTo *Trait, props map[string]interface{}, lib string) (affected int64, err error) {
	affected, err = relationTraitsInsert(eng, relation, traitFrom, traitTo, props, lib)
	return
}
