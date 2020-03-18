package lib

import (
	"time"

	"github.com/go-xorm/xorm"
	"github.com/google/uuid"
)

// Document -
type Document struct {
	ID        int64
	Dt        time.Time              `xorm:"notnull index"`
	UUID      string                 `xorm:"varchar(36) notnull unique"`
	DocTypeID int64                  `xorm:"notnull index"`
	DebID     int64                  `xorm:"notnull index"`
	CredID    int64                  `xorm:"notnull index"`
	Props     map[string]interface{} `xorm:"jsonb"`
	Status    int                    `xorm:"notnull index"`
	Created   time.Time              `xorm:"created"`
	Updated   time.Time              `xorm:"updated"`
}

// NewDocument -
func NewDocument(docTypeID int64, props map[string]interface{}) (doc *Document, err error) {
	doc = &Document{
		UUID:      uuid.New().String(),
		DocTypeID: docTypeID,
		Props:     props,
	}
	return
}

// FillDocument -
func (doc *Document) FillDocument(eng *xorm.Engine) (err error) {
	typedoc := &TypeDocument{ID: doc.DocTypeID}

	_, err = typedoc.Get(eng)
	if err != nil {
		return
	}

	err = fillDocument(doc, *typedoc)

	return
}

// Insert -
func (doc *Document) Insert(eng *xorm.Engine) (affected int64, err error) {
	affected, err = eng.Insert(doc)
	return
}

// Get -
func (doc *Document) Get(eng *xorm.Engine) (has bool, err error) {
	has, err = eng.Get(doc)
	return
}

// Update -
func (doc *Document) Update(eng *xorm.Engine) (affected int64, err error) {
	affected, err = eng.ID(doc.ID).Update(doc)
	return
}
