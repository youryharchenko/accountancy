package accountancy

import (
	"encoding/json"
	"fmt"
	"log"
)

// RunInsert -
func RunInsert(service string, request string, meta *Meta, db DB) (response string, err error) {

	log.Println("RunInsert", service, request)

	reqMap := map[string]interface{}{}
	if err = json.Unmarshal([]byte(request), &reqMap); err != nil {
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}

	bodyMap, isObj := reqMap["body"].(map[string]interface{})
	if !isObj {
		err = fmt.Errorf("RunInsert: body is not object")
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
	}

	if meta == nil {
		meta, err = LoadMeta(db)
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}
	}

	switch service {
	case "object":
		response, err = InsertOrUpdateObject(db, bodyMap, meta)
	case "link":
		response, err = InsertOrUpdateLink(db, bodyMap, meta)
	default:
		err = fmt.Errorf("RunInsert: unknown service '%s'", service)
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
	}
	if err != nil {
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}

	return
}

// RunSelect -
func RunSelect(service string, request string, meta *Meta, db DB) (response string, err error) {

	log.Println("RunSelect", service, request)

	reqMap := map[string]interface{}{}
	if err = json.Unmarshal([]byte(request), &reqMap); err != nil {
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}

	bodyMap, isObj := reqMap["body"].(map[string]interface{})
	if !isObj {
		err = fmt.Errorf("RunSelect: body is not object")
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
	}

	if meta == nil {
		meta, err = LoadMeta(db)
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}
	}

	switch service {
	case "object":
		response, err = SelectObject(db, bodyMap, meta)
	//case "link":
	//	response, err = InsertOrUpdateLink(db, objMap, meta)
	default:
		err = fmt.Errorf("RunInsert: unknown service '%s'", service)
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
	}
	if err != nil {
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}

	return
}
