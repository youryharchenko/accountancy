package lib

import (
	"time"

	"github.com/go-xorm/xorm"
	"github.com/google/uuid"
)

// Service -
type Service struct {
	ID      int64
	UUID    string                 `xorm:"varchar(36) notnull unique"`
	Name    string                 `xorm:"varchar(1024) notnull unique"`
	Props   map[string]interface{} `xorm:"jsonb"`
	Lib     string                 `xorm:"text"`
	Status  int                    `xorm:"notnull index"`
	Created time.Time              `xorm:"created"`
	Updated time.Time              `xorm:"updated"`
}

// ImportFromService -
func ImportFromService(eng *xorm.Engine, ctx string) (docs []*Document, err error) {
	docs, err = makeDocsFromContext(eng, ctx)
	return
}

// NewService -
func NewService(name string, props map[string]interface{}, lib string) (service *Service, err error) {
	service = &Service{
		UUID:  uuid.New().String(),
		Name:  name,
		Props: props,
		Lib:   lib,
	}
	return
}

// Insert -
func (service *Service) Insert(eng *xorm.Engine) (affected int64, err error) {
	affected, err = eng.Insert(service)
	return
}

// Get -
func (service *Service) Get(eng *xorm.Engine) (has bool, err error) {
	has, err = eng.Get(service)
	return
}

// Update -
func (service *Service) Update(eng *xorm.Engine) (affected int64, err error) {
	affected, err = eng.ID(service.ID).Update(service)
	return
}
