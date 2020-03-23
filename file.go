package accountancy

import (
	"time"

	"github.com/google/uuid"
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
func (file *File) Insert(db DB) (affected int64, err error) {
	affected, err = db.Insert(file)
	return
}

// Get -
func (file *File) Get(db DB) (has bool, err error) {
	has, err = db.Get(file)
	return
}

// Update -
func (file *File) Update(db DB) (affected int64, err error) {
	affected, err = db.ID(file.ID).Update(file)
	return
}

// InsertOrUpdate -
func (file *File) InsertOrUpdate(db DB, find *File) (affected int64, inserted bool, err error) {

	ok, err := find.Get(db)
	if err != nil {
		return
	}

	if ok {
		file.ID = find.ID
		_, err = file.Update(db)
		if err != nil {
			return
		}
		inserted = false
	} else {
		_, err = file.Insert(db)
		if err != nil {
			return
		}
		inserted = true
	}

	return
}
