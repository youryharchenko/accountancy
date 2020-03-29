package accountancy

import (
	"encoding/json"
	"fmt"
	"log"

	"xorm.io/xorm"
)

// RunInsert -
func RunInsert(service string, request string, meta *Meta, db DB) (response string, err error) {

	log.Println("RunInsert", service, request)

	var eng *xorm.Engine
	var sess *xorm.Session

	if db == nil {
		eng, err = ConnectRequestDB(request)
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}
		defer eng.Close()

		sess = eng.NewSession()
		defer sess.Close()

		err = sess.Begin()
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}

		db = sess
	} else {
		switch db.(type) {
		case *xorm.Engine:
			eng = db.(*xorm.Engine)
			sess = eng.NewSession()
			defer sess.Close()

			err = sess.Begin()
			if err != nil {
				response = fmt.Sprintf(tmplResponse, err.Error(), -1)
				return
			}
		case *xorm.Session:
			sess = db.(*xorm.Session)
		}
	}

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
		meta, err = LoadMeta(sess)
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}
	}

	switch service {
	case "object":
		response, err = InsertOrUpdateObject(sess, objMap, meta)
	case "link":
		//response, err = InsertOrUpdateLink(sess, request, meta)
	default:
		err = fmt.Errorf("RunInsert: unknown service '%s'", service)
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
	}
	if err != nil {
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		sess.Rollback()
		return
	}

	err = sess.Commit()
	if err != nil {
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
	}
	return
}
