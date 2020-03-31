package accountancy

import (
	"fmt"
	"log"
	"time"
)

// LinkObjects -
type LinkObjects struct {
	ID               int64                  `json:"-"`
	RelationTraitsID int64                  `xorm:"notnull unique(relation_traits_link_objects)" json:"-"`
	RelationName     string                 `xorm:"-" json:"relation"`
	TraitFromName    string                 `xorm:"-" json:"trait-from"`
	ObjectFromID     int64                  `xorm:"notnull unique(relation_traits_link_objects)" json:"-"`
	ObjectFromName   string                 `xorm:"-" json:"object-from"`
	TraitToName      string                 `xorm:"-" json:"trait-to"`
	ObjectToID       int64                  `xorm:"notnull unique(relation_traits_link_objects)" json:"-"`
	ObjectToName     string                 `xorm:"-" json:"object-to"`
	Props            map[string]interface{} `xorm:"jsonb" json:"props"`
	Hash             string                 `xorm:"text" json:"hash"`
	Status           int                    `xorm:"notnull index" json:"-"`
	Created          time.Time              `xorm:"created" json:"-"`
	Updated          time.Time              `xorm:"updated" json:"-"`
}

// Insert -
func (link *LinkObjects) Insert(db DB) (affected int64, err error) {
	affected, err = db.Insert(link)
	return
}

// Get -
func (link *LinkObjects) Get(db DB) (has bool, err error) {
	has, err = db.Get(link)
	return
}

// Update -
func (link *LinkObjects) Update(db DB) (affected int64, err error) {
	affected, err = db.ID(link.ID).Update(link)
	return
}

// InsertOrUpdate -
func (link *LinkObjects) InsertOrUpdate(db DB, find *LinkObjects) (affected int64, inserted bool, err error) {

	ok, err := find.Get(db)
	if err != nil {
		return
	}

	if ok {
		link.ID = find.ID
		_, err = link.Update(db)
		if err != nil {
			return
		}
		inserted = false
	} else {
		_, err = link.Insert(db)
		if err != nil {
			return
		}
		inserted = true
	}

	return
}

// AddLinkObjects -
func (relTraits *RelationTraits) AddLinkObjects(db DB, objectFrom *Object, objectTo *Object, props map[string]interface{}) (affected int64, inserted bool, err error) {
	affected, inserted, err = LinkObjectsInsertOrUpdate(db, relTraits, objectFrom, objectTo, props)
	return
}

// LinkObjectsInsertOrUpdate -
func LinkObjectsInsertOrUpdate(db DB, relTraits *RelationTraits, objectFrom *Object, objectTo *Object, props map[string]interface{}) (affected int64, inserted bool, err error) {

	return
}

// InsertOrUpdateLink -
func InsertOrUpdateLink(db DB, objMap map[string]interface{}, meta *Meta) (response string, err error) {
	if meta == nil {
		meta, err = LoadMeta(db)
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}
	}

	relName, ok := objMap["relation"].(string)
	if !ok {
		err = fmt.Errorf("InsertOrUpdateLink: field 'relation' is missing, source: %v", objMap)
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}
	if len(relName) == 0 {
		err = fmt.Errorf("InsertOrUpdateLink: field 'relation' length == 0, source: %v", objMap)
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}
	traitFromName, ok := objMap["trait-from"].(string)
	if !ok {
		err = fmt.Errorf("InsertOrUpdateLink: field 'trait-from' is missing, source: %v", objMap)
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}
	if len(traitFromName) == 0 {
		err = fmt.Errorf("InsertOrUpdateLink: field 'trait-from' length == 0, source: %v", objMap)
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}
	traitToName, ok := objMap["trait-to"].(string)
	if !ok {
		err = fmt.Errorf("InsertOrUpdateLink: field 'trait-to' is missing, source: %v", objMap)
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}
	if len(traitToName) == 0 {
		err = fmt.Errorf("InsertOrUpdateLink: field 'trait-to' length == 0, source: %v", objMap)
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}
	objectFromName, ok := objMap["object-from"].(string)
	if !ok {
		err = fmt.Errorf("InsertOrUpdateLink: field 'object-from' is missing, source: %v", objMap)
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}
	if len(objectFromName) == 0 {
		err = fmt.Errorf("InsertOrUpdateLink: field 'object-from' length == 0, source: %v", objMap)
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}
	objectToName, ok := objMap["object-to"].(string)
	if !ok {
		err = fmt.Errorf("InsertOrUpdateLink: field 'object-to' is missing, source: %v", objMap)
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}
	if len(objectToName) == 0 {
		err = fmt.Errorf("InsertOrUpdateLink: field 'object-to' length == 0, source: %v", objMap)
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}
	props, ok := objMap["props"].(map[string]interface{})
	if !ok {
		err = fmt.Errorf("InsertOrUpdateLink: field 'props' is missing, source: %v", objMap)
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}

	log.Println("InsertOrUpdateLink:", relName, traitFromName, objectFromName, traitToName, objectToName)

	relTraitsKey := RelationTraitsKey{
		RelationName:  relName,
		TraitFromName: traitFromName,
		TraitToName:   traitToName,
	}
	relTraits, isRelTraits := meta.RelTraitsByName[relTraitsKey]
	if !isRelTraits {
		err = fmt.Errorf("InsertOrUpdateLink: RelationTraits is missing, key: %v", relTraitsKey)
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}

	objFrom := &Object{Name: objectFromName}
	okFrom, err := objFrom.Get(db)
	if err != nil {
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}
	if !okFrom {
		err = fmt.Errorf("InsertOrUpdateLink: object '%s' is missing", objectFromName)
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}

	objTo := &Object{Name: objectToName}
	okTo, err := objTo.Get(db)
	if err != nil {
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}
	if !okTo {
		err = fmt.Errorf("InsertOrUpdateLink: object '%s' is missing", objectToName)
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}

	/*
		js := relTraits.LibContent

		vm := otto.New()
		_, err = vm.Run(js)
		if err != nil {
			return
		}

		vm.Set("makeHash", func(props map[string]interface{}) string { s, _ := makeHash(props); return s })
		vm.Set("log", func(mess string) { log.Println(mess) })

		res, err := vm.Call("constructor", nil, relTraits.Props, objFrom.Props, objTo.Props)
		if err != nil {
			return
		}
	*/

	h, err := makeHash(props)
	if err != nil {
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}

	link := &LinkObjects{
		RelationTraitsID: relTraits.ID,
		ObjectFromID:     objFrom.ID,
		ObjectToID:       objTo.ID,
		Props:            props,
		Hash:             h,
	}

	_, _, err = link.InsertOrUpdate(db, &LinkObjects{RelationTraitsID: relTraits.ID, ObjectFromID: objFrom.ID, ObjectToID: objTo.ID})
	if err != nil {
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}

	return
}
