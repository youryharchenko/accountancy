package accountancy

import (
	"encoding/json"
	"fmt"
	"log"

	//"github.com/go-xorm/xorm"
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
		response, err = ImportMeta(eng, request)
	default:
		err = fmt.Errorf("RunImport: unknown service '%s'", service)
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
	}

	return
}

// ImportMeta -
func ImportMeta(eng *xorm.Engine, request string) (response string, err error) {

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

		_, err = relation.Insert(eng)
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

		_, err = trait.Insert(eng)
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}

		log.Println(*trait)

		traitDict[trait.Name] = trait

	}

	response = fmt.Sprintf(tmplResponse, "ok", 0)

	return
}
