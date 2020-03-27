package accountancy

import (
	"time"

	"github.com/google/uuid"
)

// Relation -
type Relation struct {
	ID      int64                  `json:"-"`
	UUID    string                 `xorm:"varchar(36) notnull unique" json:"uuid"`
	Name    string                 `xorm:"varchar(1024) notnull unique" json:"name"`
	Props   map[string]interface{} `xorm:"jsonb" json:"props"`
	Status  int                    `xorm:"notnull index" json:"-"`
	Created time.Time              `xorm:"created" json:"-"`
	Updated time.Time              `xorm:"updated" json:"-"`
}

// RelationTraits -
type RelationTraits struct {
	ID            int64                  `json:"-"`
	RelationID    int64                  `xorm:"notnull unique(relation_traits)" json:"-"`
	RelationName  string                 `xorm:"-" json:"relation"`
	TraitFromID   int64                  `xorm:"notnull unique(relation_traits)" json:"-"`
	TraitFromName string                 `xorm:"-" json:"trait-from"`
	TraitToID     int64                  `xorm:"notnull unique(relation_traits)" json:"-"`
	TraitToName   string                 `xorm:"-" json:"trait-to"`
	Props         map[string]interface{} `xorm:"jsonb" json:"props"`
	Lib           string                 `xorm:"text"  json:"lib"`
	Status        int                    `xorm:"notnull index" json:"-"`
	Created       time.Time              `xorm:"created" json:"-"`
	Updated       time.Time              `xorm:"updated" json:"-"`
}

// RelationTraitsKey -
type RelationTraitsKey struct {
	RelationName  string
	TraitFromName string
	TraitToName   string
}

// LinkObjects -
type LinkObjects struct {
	ID               int64
	RelationTraitsID int64                  `xorm:"notnull unique(relation_traits_link_objects)"`
	RelationName     string                 `xorm:"-" json:"relation"`
	ObjectID         int64                  `xorm:"notnull unique(relation_traits_link_objects)"`
	ObjectName       string                 `xorm:"-" json:"object-name"`
	ObjectUUID       string                 `xorm:"-" json:"object-uuid"`
	Props            map[string]interface{} `xorm:"jsonb"`
	Status           int                    `xorm:"notnull index"`
	Created          time.Time              `xorm:"created"`
	Updated          time.Time              `xorm:"updated"`
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
func (relation *Relation) Insert(db DB) (affected int64, err error) {
	affected, err = db.Insert(relation)
	return
}

// Get -
func (relation *Relation) Get(db DB) (has bool, err error) {
	has, err = db.Get(relation)
	return
}

// Update -
func (relation *Relation) Update(db DB) (affected int64, err error) {
	affected, err = db.ID(relation.ID).Update(relation)
	return
}

// InsertOrUpdate -
func (relation *Relation) InsertOrUpdate(db DB, find *Relation) (affected int64, inserted bool, err error) {

	ok, err := find.Get(db)
	if err != nil {
		return
	}

	if ok {
		relation.ID = find.ID
		_, err = relation.Update(db)
		if err != nil {
			return
		}
		inserted = false
	} else {
		_, err = relation.Insert(db)
		if err != nil {
			return
		}
		inserted = true
	}

	return
}

// AddRelationTraits -
func (relation *Relation) AddRelationTraits(db DB, traitFrom *Trait, traitTo *Trait, props map[string]interface{}, lib string) (affected int64, inserted bool, err error) {
	affected, inserted, err = relationTraitsInsertOrUpdate(db, relation, traitFrom, traitTo, props, lib)
	return
}
