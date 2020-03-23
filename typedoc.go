package accountancy

import (
	"time"

	"github.com/google/uuid"
)

// TypeDocument -
type TypeDocument struct {
	ID            int64
	UUID          string                 `xorm:"varchar(36) notnull unique"`
	Name          string                 `xorm:"varchar(1024) notnull unique"`
	DebTraitID    int64                  `xorm:"notnull index"`
	DebTraitName  string                 `xorm:"-" json:"deb-trait"`
	CredTraitID   int64                  `xorm:"notnull index"`
	CredTraitName string                 `xorm:"-" json:"cred-trait"`
	Props         map[string]interface{} `xorm:"jsonb"`
	Lib           string                 `xorm:"text"`
	Status        int                    `xorm:"notnull index"`
	Created       time.Time              `xorm:"created"`
	Updated       time.Time              `xorm:"updated"`
}

// NewTypeDocument -
func NewTypeDocument(name string, debTraitID int64, credTraitID int64, props map[string]interface{}, lib string) (typedoc *TypeDocument, err error) {
	typedoc = &TypeDocument{
		UUID:        uuid.New().String(),
		Name:        name,
		DebTraitID:  debTraitID,
		CredTraitID: credTraitID,
		Props:       props,
		Lib:         lib,
	}
	return
}

// Insert -
func (typedoc *TypeDocument) Insert(db DB) (affected int64, err error) {
	affected, err = db.Insert(typedoc)
	return
}

// Get -
func (typedoc *TypeDocument) Get(db DB) (has bool, err error) {
	has, err = db.Get(typedoc)
	return
}

// Update -
func (typedoc *TypeDocument) Update(db DB) (affected int64, err error) {
	affected, err = db.ID(typedoc.ID).Update(typedoc)
	return
}

// InsertOrUpdate -
func (typedoc *TypeDocument) InsertOrUpdate(db DB, find *TypeDocument) (affected int64, inserted bool, err error) {

	ok, err := find.Get(db)
	if err != nil {
		return
	}

	if ok {
		typedoc.ID = find.ID
		_, err = typedoc.Update(db)
		if err != nil {
			return
		}
		inserted = false
	} else {
		_, err = typedoc.Insert(db)
		if err != nil {
			return
		}
		inserted = true
	}

	return
}

/*
// SessionInsert -
func (typedoc *TypeDocument) SessionInsert(sess *xorm.Session) (affected int64, err error) {
	affected, err = sess.Insert(typedoc)
	return
}

// SessionGet -
func (typedoc *TypeDocument) SessionGet(sess *xorm.Session) (has bool, err error) {
	has, err = sess.Get(typedoc)
	return
}

// SessionUpdate -
func (typedoc *TypeDocument) SessionUpdate(sess *xorm.Session) (affected int64, err error) {
	affected, err = sess.ID(typedoc.ID).Update(typedoc)
	return
}
*/
