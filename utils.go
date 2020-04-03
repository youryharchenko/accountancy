package accountancy

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/robertkrimen/otto"
	"github.com/tidwall/gjson"
	"xorm.io/xorm"
)

func testProps(vm *otto.Otto, props interface{}, filterProps string) (ok bool, err error) {
	vmc := vm.Copy()

	vmc.Set("props", props)
	var res otto.Value
	res, err = vmc.Eval(filterProps)
	if err != nil {
		return
	}
	ok, err = res.ToBoolean()
	return
}

func makeFields(vm *otto.Otto, ent interface{}, fields string) (final interface{}, err error) {
	if fields == "*" {
		final = ent
		return
	}
	fields = "var r = " + strings.ReplaceAll(fields, "'", "\"") + "; r"
	vmc := vm.Copy()
	vmc.Set("obj", ent)
	res, err := vmc.Eval(fields)
	if err != nil {
		return
	}
	final, err = res.Export()
	return
}

func relationTraitsInsert(db DB, relation *Relation, traitFrom *Trait, traitTo *Trait, props map[string]interface{}, lib string) (affected int64, err error) {
	relationTraits := &RelationTraits{
		RelationID:  relation.ID,
		TraitFromID: traitFrom.ID,
		TraitToID:   traitTo.ID,
		Props:       props,
		Lib:         lib,
	}
	affected, err = db.Insert(relationTraits)
	return
}

func relationTraitsInsertOrUpdate(db DB, relation *Relation, traitFrom *Trait, traitTo *Trait, props map[string]interface{}, lib string) (affected int64, inserted bool, err error) {
	relationTraits := &RelationTraits{
		RelationID:  relation.ID,
		TraitFromID: traitFrom.ID,
		TraitToID:   traitTo.ID,
		Props:       props,
		Lib:         lib,
	}

	find := &RelationTraits{
		RelationID:  relation.ID,
		TraitFromID: traitFrom.ID,
		TraitToID:   traitTo.ID,
	}

	ok, err := db.Get(find)
	if err != nil {
		return
	}

	if ok {
		relationTraits.ID = find.ID
		_, err = db.ID(find.ID).Update(relationTraits)
		if err != nil {
			return
		}
		inserted = false
	} else {
		_, err = db.Insert(relationTraits)
		if err != nil {
			return
		}
		inserted = true
	}

	return
}

func traitObjectUpdate(db DB, trait *Trait, object *Object, props map[string]interface{}) (affected int64, err error) {
	traitObject := &TraitObject{
		TraitID:  trait.ID,
		ObjectID: object.ID,
	}

	_, err = db.Get(traitObject)
	if err != nil {
		return
	}

	traitObject.Props = props

	affected, err = db.ID(traitObject.ID).Update(traitObject)
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

func makeHash(props map[string]interface{}) (s string, err error) {
	buf, err := json.Marshal(props)
	if err != nil {
		return
	}
	h := sha256.New()
	_, err = h.Write(buf)
	if err != nil {
		return
	}
	s = base64.StdEncoding.EncodeToString(h.Sum(nil))
	return
}
