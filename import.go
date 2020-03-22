package accountancy

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/tidwall/gjson"
	"xorm.io/xorm"
)

// RunImport -
func RunImport(service string, request string) (response string, err error) {

	var eng *xorm.Engine
	eng, err = ConnectRequestDB(request)
	if err != nil {
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}
	defer eng.Close()

	switch service {
	case "meta":
		sess := eng.NewSession()
		defer sess.Close()

		err = sess.Begin()
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}

		response, err = ImportMeta(sess, request)
		if err != nil {
			err = sess.Rollback()
			if err != nil {
				response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			}
			return
		}

		err = sess.Commit()
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		}

	default:
		err = fmt.Errorf("RunImport: unknown service '%s'", service)
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
	}

	return
}

// ImportMeta -
func ImportMeta(sess *xorm.Session, request string) (response string, err error) {

	relDict := map[string]*Relation{}
	traitDict := map[string]*Trait{}

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
	}
	for _, relJSON := range relations {
		relation := &Relation{}
		err = json.Unmarshal([]byte(relJSON.Raw), relation)
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}

		_, err = relation.SessionInsert(sess)
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}

		log.Println(*relation)

		relDict[relation.Name] = relation

	}

	var traits []gjson.Result
	traitsJSON := bodyJSON.Get("traits")
	if traitsJSON.Exists() {
		traits = traitsJSON.Array()
	}
	for _, traitJSON := range traits {
		trait := &Trait{}
		err = json.Unmarshal([]byte(traitJSON.Raw), trait)
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}

		_, err = trait.SessionInsert(sess)
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}

		log.Println(*trait)

		traitDict[trait.Name] = trait

	}

	var relTraits []gjson.Result
	relTraitsJSON := bodyJSON.Get("relation-traits")
	if relTraitsJSON.Exists() {
		relTraits = relTraitsJSON.Array()
	}
	for _, relTraitJSON := range relTraits {
		relationTraits := &RelationTraits{}
		err = json.Unmarshal([]byte(relTraitJSON.Raw), relationTraits)
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}

		relation, ok := relDict[relationTraits.RelationName]
		if !ok {
			err = fmt.Errorf("ImportMeta: relation '%s' is undefined", relationTraits.RelationName)
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}

		traitFrom, ok := traitDict[relationTraits.TraitFromName]
		if !ok {
			err = fmt.Errorf("ImportMeta: trait 'from' '%s' is undefined", relationTraits.TraitFromName)
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}

		traitTo, ok := traitDict[relationTraits.TraitToName]
		if !ok {
			err = fmt.Errorf("ImportMeta: trait 'to' '%s' is undefined", relationTraits.TraitToName)
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}

		_, err = relation.SessionAddRelationTraits(sess, traitFrom, traitTo, relationTraits.Props, relationTraits.Lib)
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}

	}

	var typeDocs []gjson.Result
	typeDocsJSON := bodyJSON.Get("typedocs")
	if typeDocsJSON.Exists() {
		typeDocs = typeDocsJSON.Array()
	}
	for _, typeDocJSON := range typeDocs {
		typeDoc := &TypeDocument{}
		err = json.Unmarshal([]byte(typeDocJSON.Raw), typeDoc)
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}

		debTrait, ok := traitDict[typeDoc.DebTraitName]
		if !ok {
			err = fmt.Errorf("ImportMeta: trait 'deb' '%s' is undefined", typeDoc.DebTraitName)
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}

		credTrait, ok := traitDict[typeDoc.CredTraitName]
		if !ok {
			err = fmt.Errorf("ImportMeta: trait 'cred' '%s' is undefined", typeDoc.CredTraitName)
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}

		typeDoc.DebTraitID = debTrait.ID
		typeDoc.CredTraitID = credTrait.ID

		_, err = typeDoc.SessionInsert(sess)
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}

	}

	response = fmt.Sprintf(tmplResponse, "ok", 0)

	return
}
