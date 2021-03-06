package accountancy

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/robertkrimen/otto"
	"xorm.io/xorm"
)

// Object -
type Object struct {
	ID      int64                  `json:"-"`
	UUID    string                 `xorm:"varchar(36) notnull unique" json:"uuid"`
	Name    string                 `xorm:"varchar(1024) notnull unique" json:"name"`
	Props   map[string]interface{} `xorm:"jsonb" json:"props"`
	Hash    string                 `xorm:"text" json:"hash"`
	Status  int                    `xorm:"notnull index" json:"status"`
	Created time.Time              `xorm:"created" json:"-"`
	Updated time.Time              `xorm:"updated" json:"-"`
}

// NewObject -
func NewObject(name string, props map[string]interface{}) (object *Object, err error) {
	h, err := makeHash(props)
	if err != nil {
		return
	}
	object = &Object{
		UUID:  uuid.New().String(),
		Name:  name,
		Props: props,
		Hash:  h,
	}
	return
}

// Insert -
func (object *Object) Insert(db DB) (affected int64, err error) {
	affected, err = db.Insert(object)
	return
}

// Get -
func (object *Object) Get(db DB) (has bool, err error) {
	has, err = db.Get(object)
	return
}

// Update -
func (object *Object) Update(db DB) (affected int64, err error) {
	affected, err = db.ID(object.ID).Update(object)
	return
}

// InsertOrUpdate -
func (object *Object) InsertOrUpdate(db DB, find *Object) (affected int64, inserted bool, err error) {

	ok, err := find.Get(db)
	if err != nil {
		return
	}

	if ok {
		object.ID = find.ID
		_, err = object.Update(db)
		if err != nil {
			return
		}
		inserted = false
	} else {
		_, err = object.Insert(db)
		if err != nil {
			return
		}
		inserted = true
	}

	return
}

// AddTrait -
func (object *Object) AddTrait(eng *xorm.Engine, trait *Trait, meta *Meta) (affected int64, err error) {
	affected, _, err = TraitObjectInsertOrUpdate(eng, trait, object, meta)
	return
}

// InsertOrUpdateObject -
func InsertOrUpdateObject(db DB, objMap map[string]interface{}, meta *Meta) (response string, err error) {
	if meta == nil {
		meta, err = LoadMeta(db)
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}
	}

	name, ok := objMap["name"].(string)
	if !ok {
		err = fmt.Errorf("InsertOrUpdateObject: field 'name' is missing, source: %v", objMap)
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}
	if len(name) == 0 {
		err = fmt.Errorf("InsertOrUpdateObject: field 'name' length == 0, source: %v", objMap)
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}
	props, ok := objMap["props"].(map[string]interface{})
	if !ok {
		err = fmt.Errorf("InsertOrUpdateObject: field 'props' is missing, source: %v", objMap)
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}
	traits, ok := objMap["traits"].([]interface{})
	if !ok {
		err = fmt.Errorf("InsertOrUpdateObject: field 'traits' is missing, source: %v", objMap)
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}

	log.Println("InsertOrUpdateObject:", name)

	var obj *Object
	obj, err = NewObject(name, props)
	if err != nil {
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}

	_, _, err = obj.InsertOrUpdate(db, &Object{Name: name})
	if err != nil {
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}

	for _, traitName := range traits {

		name, ok := traitName.(string)
		if !ok {
			err = fmt.Errorf("InsertOrUpdateObject: trait name '%s' is not string", traitName)
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}
		trait, ok := meta.TraitsByName[name]
		if !ok {
			err = fmt.Errorf("InsertOrUpdateObject: trait '%s' is missing", traitName)
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}
		_, _, err = trait.AddObject(db, obj, meta)
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}
	}

	return
}

// SelectObject -
func SelectObject(db DB, objMap map[string]interface{}, meta *Meta) (response string, err error) {
	if meta == nil {
		meta, err = LoadMeta(db)
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}
	}

	filterMap, ok := objMap["filter"].(map[string]interface{})
	if !ok {
		filterMap = map[string]interface{}{}
	}

	condition, ok := filterMap["condition"].(string)
	if !ok {
		condition = ""
	}
	params, ok := filterMap["params"].([]interface{})
	if !ok {
		params = []interface{}{}
	}

	filterProps, ok := objMap["filterProps"].(string)
	if !ok {
		filterProps = "true"
	}

	orderBy, ok := objMap["orderBy"].(string)
	if !ok {
		orderBy = ""
	}

	skip, ok := objMap["skip"].(float64)
	if !ok {
		skip = float64(0)
	}

	limit, ok := objMap["limit"].(float64)
	if !ok {
		limit = float64(10)
	}

	fields, ok := objMap["fields"].(string)
	if !ok {
		fields = "*"
	}

	log.Println("SelectObject:", filterMap, skip, limit)

	list := []Object{}
	err = db.Where(condition, params...).
		OrderBy(orderBy).
		//Limit(int(limit), int(skip)).
		Find(&list)
	if err != nil {
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}

	log.Println(len(list))
	result := []interface{}{}

	vm := otto.New()
	count := int64(0)
	for j, obj := range list {
		if int64(limit) <= count {
			break
		}
		if int64(skip) > int64(j) {
			continue
		}

		var b bool
		b, err = testProps(vm, obj.Props, filterProps)
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}
		if b {
			count++
			var fieldsMap interface{}
			fieldsMap, err = makeFields(vm, obj, fields)
			if err != nil {
				response = fmt.Sprintf(tmplResponse, err.Error(), -1)
				return
			}
			result = append(result, fieldsMap)
		}
	}

	respMap := map[string]interface{}{
		"message": "ok",
		"status":  0,
		"body": map[string]interface{}{
			"filter":      filterMap,
			"skip":        skip,
			"limit":       limit,
			"filterProps": filterProps,
			"count":       count,
			"fields":      fields,
			"result":      result,
		},
	}

	respBuf, err := json.MarshalIndent(respMap, "", "  ")
	if err != nil {
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}

	response = string(respBuf)

	return
}
