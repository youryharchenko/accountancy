package accountancy

import (
	"time"

	//"github.com/go-xorm/xorm"
	"github.com/google/uuid"
	"xorm.io/xorm"
)

// TypeDocument -
type TypeDocument struct {
	ID          int64
	UUID        string                 `xorm:"varchar(36) notnull unique"`
	Name        string                 `xorm:"varchar(1024) notnull unique"`
	DebTraitID  int64                  `xorm:"notnull index"`
	CredTraitID int64                  `xorm:"notnull index"`
	Props       map[string]interface{} `xorm:"jsonb"`
	Lib         string                 `xorm:"text"`
	Status      int                    `xorm:"notnull index"`
	Created     time.Time              `xorm:"created"`
	Updated     time.Time              `xorm:"updated"`
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
func (typedoc *TypeDocument) Insert(eng *xorm.Engine) (affected int64, err error) {
	affected, err = eng.Insert(typedoc)
	return
}

// Get -
func (typedoc *TypeDocument) Get(eng *xorm.Engine) (has bool, err error) {
	has, err = eng.Get(typedoc)
	return
}

// Update -
func (typedoc *TypeDocument) Update(eng *xorm.Engine) (affected int64, err error) {
	affected, err = eng.ID(typedoc.ID).Update(typedoc)
	return
}
