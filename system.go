package accountancy

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
)

// RunSystem -
func RunSystem(service string, request string, meta *Meta, db DB) (response string, err error) {

	switch service {
	case "initdb":
		err = SyncAll(db)
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		} else {
			response = fmt.Sprintf(tmplResponse, "ok", 0)
		}
	case "dropdb":
		err = DropAll(db)
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		} else {
			response = fmt.Sprintf(tmplResponse, "ok", 0)
		}
	case "upload":
		response, err = SystemUpload(db, request)
	default:
		err = fmt.Errorf("RunSystem: unknown service '%s'", service)
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
	}

	return
}

// SystemUpload -
func SystemUpload(db DB, request string) (response string, err error) {

	srcJSON := gjson.Parse(request)

	bodyJSON := srcJSON.Get("body")
	if !bodyJSON.Exists() {
		err = fmt.Errorf("SystemUpload: object 'body' is missing")
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}

	uuid := bodyJSON.Get("uuid").String()
	name := bodyJSON.Get("name").String()
	props := map[string]interface{}{}
	err = json.Unmarshal([]byte(bodyJSON.Get("props").Raw), &props)
	if err != nil {
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}

	content, err := base64.StdEncoding.DecodeString(bodyJSON.Get("content").String())
	if err != nil {
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}

	file := &File{
		UUID:    uuid,
		Name:    name,
		Props:   props,
		Content: content,
	}

	_, inserted, err := file.InsertOrUpdate(db, &File{Name: name})

	mess := "updated"
	if inserted {
		mess = "inserted"
	}
	response = fmt.Sprintf(tmplResponse, mess, 0)

	return
}
