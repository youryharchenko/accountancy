package accountancy

import (
	"encoding/json"
	"fmt"
	"log"

	//"github.com/go-xorm/xorm"
	"github.com/robertkrimen/otto"
	"github.com/tidwall/gjson"
	"xorm.io/xorm"
)

func relationTraitsInsert(eng *xorm.Engine, relation *Relation, traitFrom *Trait, traitTo *Trait, props map[string]interface{}, lib string) (affected int64, err error) {
	relationTraits := &RelationTraits{
		RelationID:  relation.ID,
		TraitFromID: traitFrom.ID,
		TraitToID:   traitTo.ID,
		Props:       props,
		Lib:         lib,
	}
	affected, err = eng.Insert(relationTraits)
	return
}

func traitObjectInsert(eng *xorm.Engine, trait *Trait, object *Object) (affected int64, err error) {
	traitObject := &TraitObject{
		TraitID:  trait.ID,
		ObjectID: object.ID,
		Props:    trait.Props,
	}
	affected, err = eng.Insert(traitObject)
	return
}

func traitObjectUpdate(eng *xorm.Engine, trait *Trait, object *Object, props map[string]interface{}) (affected int64, err error) {
	traitObject := &TraitObject{
		TraitID:  trait.ID,
		ObjectID: object.ID,
	}

	_, err = eng.Get(traitObject)
	if err != nil {
		return
	}

	traitObject.Props = props

	affected, err = eng.ID(traitObject.ID).Update(traitObject)
	return
}

func fillDocument(doc *Document, typedoc TypeDocument) (err error) {
	var jsonDoc []byte
	var jsonTypedoc []byte
	var result otto.Value

	vm := otto.New()

	_, err = vm.Run(typedoc.Lib)
	if err != nil {
		return
	}

	jsonDoc, err = json.Marshal(*doc)
	jsonTypedoc, err = json.Marshal(typedoc)

	result, err = vm.Call("fillDocument", nil, string(jsonDoc), string(jsonTypedoc))
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(result.String()), doc)

	return

}

func makeDocsFromContext(eng *xorm.Engine, ctx string) (docs []*Document, err error) {
	var result otto.Value

	name := gjson.Parse(ctx).Get("pay.request.service").String()
	if len(name) == 0 {
		return nil, fmt.Errorf("'pay.request.service' not found in the context")
	}

	service := &Service{Name: name}

	_, err = service.Get(eng)
	if err != nil {
		return
	}
	log.Println(name, *service)

	vm := otto.New()

	_, err = vm.Run(service.Lib)
	if err != nil {
		return
	}

	result, err = vm.Call("importDocuments", nil, ctx)
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(result.String()), &docs)

	return
}
