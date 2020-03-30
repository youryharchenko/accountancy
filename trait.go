package accountancy

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/robertkrimen/otto"
)

// Trait -
type Trait struct {
	ID         int64                  `json:"-"`
	UUID       string                 `xorm:"varchar(36) notnull unique" json:"uuid"`
	Name       string                 `xorm:"varchar(1024) notnull unique" json:"name"`
	Props      map[string]interface{} `xorm:"jsonb" json:"props"`
	Lib        string                 `xorm:"text" json:"lib"`
	LibContent string                 `xorm:"-" json:"-"`
	Status     int                    `xorm:"notnull index" json:"-"`
	Created    time.Time              `xorm:"created" json:"-"`
	Updated    time.Time              `xorm:"updated" json:"-"`
}

// TraitObject -
type TraitObject struct {
	ID         int64                  `json:"-"`
	TraitID    int64                  `xorm:"notnull unique(trait_object)" json:"-"`
	TraitName  string                 `xorm:"-" json:"trait-name"`
	ObjectID   int64                  `xorm:"notnull unique(trait_object)" json:"-"`
	ObjectName string                 `xorm:"-" json:"object-name"`
	Props      map[string]interface{} `xorm:"jsonb" json:"props"`
	Hash       string                 `xorm:"text" json:"hash"`
	Status     int                    `xorm:"notnull index" json:"-"`
	Created    time.Time              `xorm:"created" json:"-"`
	Updated    time.Time              `xorm:"updated" json:"-"`
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
func (trait *Trait) Insert(db DB) (affected int64, err error) {
	affected, err = db.Insert(trait)
	return
}

// Get -
func (trait *Trait) Get(db DB) (has bool, err error) {
	has, err = db.Get(trait)
	return
}

// Update -
func (trait *Trait) Update(db DB) (affected int64, err error) {
	affected, err = db.ID(trait.ID).Update(trait)
	return
}

// InsertOrUpdate -
func (trait *Trait) InsertOrUpdate(db DB, find *Trait) (affected int64, inserted bool, err error) {

	ok, err := find.Get(db)
	if err != nil {
		return
	}

	if ok {
		trait.ID = find.ID
		_, err = trait.Update(db)
		if err != nil {
			return
		}
		inserted = false
	} else {
		_, err = trait.Insert(db)
		if err != nil {
			return
		}
		inserted = true
	}

	return
}

// AddRelationFrom -
func (trait *Trait) AddRelationFrom(db DB, relation *Relation, traitFrom *Trait, props map[string]interface{}, lib string) (affected int64, err error) {
	affected, _, err = relationTraitsInsertOrUpdate(db, relation, traitFrom, trait, props, lib)
	return
}

// AddRelationTo -
func (trait *Trait) AddRelationTo(db DB, relation *Relation, traitTo *Trait, props map[string]interface{}, lib string) (affected int64, err error) {
	affected, _, err = relationTraitsInsertOrUpdate(db, relation, trait, traitTo, props, lib)
	return
}

// AddObject -
func (trait *Trait) AddObject(db DB, object *Object, meta *Meta) (affected int64, inserted bool, err error) {
	affected, inserted, err = TraitObjectInsertOrUpdate(db, trait, object, meta)
	return
}

// UpdateObjectProps -
func (trait *Trait) UpdateObjectProps(db DB, object *Object, props map[string]interface{}) (affected int64, err error) {
	affected, err = traitObjectUpdate(db, trait, object, props)
	return
}

// TraitObjectInsertOrUpdate -
func TraitObjectInsertOrUpdate(db DB, trait *Trait, object *Object, meta *Meta) (affected int64, inserted bool, err error) {

	log.Println("TraitObjectInsertOrUpdate:", trait.Name, object.Name)

	traitObject := &TraitObject{
		TraitID:  trait.ID,
		ObjectID: object.ID,
		Props:    trait.Props,
	}

	find := &TraitObject{
		TraitID:  trait.ID,
		ObjectID: object.ID,
	}

	ok, err := db.Get(find)
	if err != nil {
		return
	}

	js := trait.LibContent

	vm := otto.New()
	_, err = vm.Run(js)
	if err != nil {
		return
	}

	vm.Set("makeHash", func(props map[string]interface{}) string { s, _ := makeHash(props); return s })
	vm.Set("log", func(mess string) { log.Println(mess) })

	res, err := vm.Call("constructor", nil, trait.Props, object.Props)
	if err != nil {
		return
	}

	result, err := res.Export()
	if err != nil {
		return
	}

	pack, isPack := result.(map[string]interface{})
	if !isPack {
		err = fmt.Errorf("TraitObjectInsertOrUpdate: result is not object")
		return
	}

	//log.Println(reflect.TypeOf(pack["status"]).String())
	status, isStatus := pack["status"].(float64)
	if !isStatus {
		err = fmt.Errorf("TraitObjectInsertOrUpdate: status is not number")
		return
	}

	if int64(status) != 0 {
		return
	}

	props, isProps := pack["props"].(map[string]interface{})
	if !isProps {
		err = fmt.Errorf("TraitObjectInsertOrUpdate: props is not object")
		return
	}

	h, isH := pack["hash"].(string)
	if !isH || len(h) == 0 {
		err = fmt.Errorf("TraitObjectInsertOrUpdate: hash is not string or is empty")
		return
	}

	traitObject.Props = props
	traitObject.Hash = h
	if err != nil {
		return
	}

	if ok {
		traitObject.ID = find.ID
		_, err = db.ID(find.ID).Update(traitObject)
		if err != nil {
			return
		}
		inserted = false
	} else {
		_, err = db.Insert(traitObject)
		if err != nil {
			return
		}
		inserted = true
	}

	objBatch, isObjBatch := pack["batch"].(map[string]interface{})
	if isObjBatch {

		var buf []byte
		buf, err = json.Marshal(objBatch)
		if err != nil {
			return
		}
		_, err = RunBatch(string(buf), meta, db)
		if err != nil {
			return
		}

	}

	return
}
