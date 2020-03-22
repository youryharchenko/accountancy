package accountancy

import (
	"time"

	"github.com/google/uuid"
	"xorm.io/xorm"
)

// File -
type File struct {
	ID      int64
	UUID    string                 `xorm:"varchar(36) notnull unique"`
	Name    string                 `xorm:"varchar(1024) notnull unique"`
	Props   map[string]interface{} `xorm:"jsonb"`
	Content []byte                 `xorm:"bytea"`
	Status  int                    `xorm:"notnull index"`
	Created time.Time              `xorm:"created"`
	Updated time.Time              `xorm:"updated"`
}

// NewFile -
func NewFile(name string, props map[string]interface{}, content []byte) (file *File, err error) {
	file = &File{
		UUID:    uuid.New().String(),
		Name:    name,
		Props:   props,
		Content: content,
	}
	return
}

// Insert -
func (file *File) Insert(eng *xorm.Engine) (affected int64, err error) {
	affected, err = eng.Insert(file)
	return
}

// Get -
func (file *File) Get(eng *xorm.Engine) (has bool, err error) {
	has, err = eng.Get(file)
	return
}

// Update -
func (file *File) Update(eng *xorm.Engine) (affected int64, err error) {
	affected, err = eng.ID(file.ID).Update(file)
	return
}

// SessionInsert -
func (file *File) SessionInsert(sess *xorm.Session) (affected int64, err error) {
	affected, err = sess.Insert(file)
	return
}

// SessionGet -
func (file *File) SessionGet(sess *xorm.Session) (has bool, err error) {
	has, err = sess.Get(file)
	return
}

// SessionUpdate -
func (file *File) SessionUpdate(sess *xorm.Session) (affected int64, err error) {
	affected, err = sess.ID(file.ID).Update(file)
	return
}
