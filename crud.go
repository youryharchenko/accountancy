package accountancy

import (
	"encoding/json"
	"fmt"
	"log"
)

// RunInsert -
func RunInsert(service string, request string, meta *Meta, db DB) (response string, err error) {

	log.Println("RunInsert", service, request)

	bodyMap := map[string]interface{}{}
	if err = json.Unmarshal([]byte(request), &bodyMap); err != nil {
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}

	objMap, isObj := bodyMap["body"].(map[string]interface{})
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
		response, err = InsertOrUpdateObject(db, objMap, meta)
	case "link":
		//response, err = InsertOrUpdateLink(sess, request, meta)
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
