package accountancy

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/tidwall/gjson"
)

// RunImport -
func RunImport(service string, request string, meta *Meta, db DB) (response string, err error) {

	log.Println("Start Import:", service)

	if meta == nil {
		meta, err = LoadMeta(db)
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}
	}

	switch service {
	case "meta":
		response, err = ImportMeta(db, request, meta)
		if err != nil {
			return
		}

		meta, err = LoadMeta(db)
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}
	case "data":
		response, err = ImportData(db, request, meta)
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}

	default:
		err = fmt.Errorf("RunImport: unknown service '%s'", service)
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
	}

	log.Println("Finish Import:", service)
	return
}

// ImportMeta -
func ImportMeta(db DB, request string, meta *Meta) (response string, err error) {

	if meta == nil {
		meta, err = LoadMeta(db)
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}
	}

	srcJSON := gjson.Parse(request)

	bodyJSON := srcJSON.Get("body")
	if !bodyJSON.Exists() {
		err = fmt.Errorf("ImportMeta: object 'body' is missing")
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}

	var relations []gjson.Result
	relsJSON := bodyJSON.Get("relations")
	if relsJSON.Exists() {
		relations = relsJSON.Array()

		for _, relJSON := range relations {
			relation := &Relation{}
			err = json.Unmarshal([]byte(relJSON.Raw), relation)
			if err != nil {
				response = fmt.Sprintf(tmplResponse, err.Error(), -1)
				return
			}

			_, _, err = relation.InsertOrUpdate(db, &Relation{Name: relation.Name})
			if err != nil {
				response = fmt.Sprintf(tmplResponse, err.Error(), -1)
				return
			}

			//log.Println(*relation)

			//relDict[relation.Name] = relation
			meta.AddRelation(relation)

		}
	}

	var traits []gjson.Result
	traitsJSON := bodyJSON.Get("traits")
	if traitsJSON.Exists() {
		traits = traitsJSON.Array()

		for _, traitJSON := range traits {
			trait := &Trait{}
			err = json.Unmarshal([]byte(traitJSON.Raw), trait)
			if err != nil {
				response = fmt.Sprintf(tmplResponse, err.Error(), -1)
				return
			}

			_, _, err = trait.InsertOrUpdate(db, &Trait{Name: trait.Name})
			if err != nil {
				response = fmt.Sprintf(tmplResponse, err.Error(), -1)
				return
			}

			//log.Println(*trait)

			//traitDict[trait.Name] = trait
			meta.AddTrait(trait)

		}
	}

	var relTraits []gjson.Result
	relTraitsJSON := bodyJSON.Get("relation-traits")
	if relTraitsJSON.Exists() {
		relTraits = relTraitsJSON.Array()

		for _, relTraitJSON := range relTraits {
			relationTraits := &RelationTraits{}
			err = json.Unmarshal([]byte(relTraitJSON.Raw), relationTraits)
			if err != nil {
				response = fmt.Sprintf(tmplResponse, err.Error(), -1)
				return
			}

			relation, ok := meta.RelationsByName[relationTraits.RelationName]
			if !ok {
				err = fmt.Errorf("ImportMeta: relation '%s' is undefined", relationTraits.RelationName)
				response = fmt.Sprintf(tmplResponse, err.Error(), -1)
				return
			}

			traitFrom, ok := meta.TraitsByName[relationTraits.TraitFromName]
			if !ok {
				err = fmt.Errorf("ImportMeta: trait 'from' '%s' is undefined", relationTraits.TraitFromName)
				response = fmt.Sprintf(tmplResponse, err.Error(), -1)
				return
			}

			traitTo, ok := meta.TraitsByName[relationTraits.TraitToName]
			if !ok {
				err = fmt.Errorf("ImportMeta: trait 'to' '%s' is undefined", relationTraits.TraitToName)
				response = fmt.Sprintf(tmplResponse, err.Error(), -1)
				return
			}

			_, _, err = relation.AddRelationTraits(db, traitFrom, traitTo, relationTraits.Props, relationTraits.Lib)
			if err != nil {
				response = fmt.Sprintf(tmplResponse, err.Error(), -1)
				return
			}

		}
	}

	var typeDocs []gjson.Result
	typeDocsJSON := bodyJSON.Get("typedocs")
	if typeDocsJSON.Exists() {
		typeDocs = typeDocsJSON.Array()

		for _, typeDocJSON := range typeDocs {
			typeDoc := &TypeDocument{}
			err = json.Unmarshal([]byte(typeDocJSON.Raw), typeDoc)
			if err != nil {
				response = fmt.Sprintf(tmplResponse, err.Error(), -1)
				return
			}

			debTrait, ok := meta.TraitsByName[typeDoc.DebTraitName]
			if !ok {
				err = fmt.Errorf("ImportMeta: trait 'deb' '%s' is undefined", typeDoc.DebTraitName)
				response = fmt.Sprintf(tmplResponse, err.Error(), -1)
				return
			}

			credTrait, ok := meta.TraitsByName[typeDoc.CredTraitName]
			if !ok {
				err = fmt.Errorf("ImportMeta: trait 'cred' '%s' is undefined", typeDoc.CredTraitName)
				response = fmt.Sprintf(tmplResponse, err.Error(), -1)
				return
			}

			typeDoc.DebTraitID = debTrait.ID
			typeDoc.CredTraitID = credTrait.ID

			_, _, err = typeDoc.InsertOrUpdate(db, &TypeDocument{Name: typeDoc.Name})
			if err != nil {
				response = fmt.Sprintf(tmplResponse, err.Error(), -1)
				return
			}

		}
	}
	response = fmt.Sprintf(tmplResponse, "ok", 0)

	return
}

// ImportData -
func ImportData(db DB, request string, meta *Meta) (response string, err error) {

	if meta == nil {
		meta, err = LoadMeta(db)
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}
	}

	srcJSON := gjson.Parse(request)

	bodyJSON := srcJSON.Get("body")
	if !bodyJSON.Exists() {
		err = fmt.Errorf("ImportData: object 'body' is missing")
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}

	var objects []gjson.Result
	objsJSON := bodyJSON.Get("objects")
	if objsJSON.Exists() {
		objects = objsJSON.Array()

		for _, objJSON := range objects {

			objMap := map[string]interface{}{}
			if err = json.Unmarshal([]byte(objJSON.Raw), &objMap); err != nil {
				response = fmt.Sprintf(tmplResponse, err.Error(), -1)
				return
			}

			response, err = InsertOrUpdateObject(db, objMap, meta)
			if err != nil {
				response = fmt.Sprintf(tmplResponse, err.Error(), -1)
				return
			}

		}
	}

	response = fmt.Sprintf(tmplResponse, "ok", 0)

	return

}
