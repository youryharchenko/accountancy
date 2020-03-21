package accountancy

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
	"xorm.io/xorm"
)

// RunSystem -
func RunSystem(service string, request string) (response string, err error) {

	var eng *xorm.Engine
	eng, err = ConnectRequestDB(request)
	if err != nil {
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}
	defer eng.Close()

	switch service {
	case "initdb":
		err = SyncAll(eng)
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		} else {
			response = fmt.Sprintf(tmplResponse, "ok", 0)
		}
	case "dropdb":
		err = DropAll(eng)
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		} else {
			response = fmt.Sprintf(tmplResponse, "ok", 0)
		}
	case "upload":
		response, err = SystemUpload(eng, request)
	default:
		err = fmt.Errorf("RunSystem: unknown service '%s'", service)
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
	}

	return
}

// SystemUpload -
func SystemUpload(eng *xorm.Engine, request string) (response string, err error) {

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

	find := &File{
		Name: name,
	}

	ok, err := find.Get(eng)
	if err != nil {
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}

	if ok {
		find.Content = content
		find.Props = props
		_, err = find.Update(eng)
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}
	} else {
		file := &File{
			UUID:    uuid,
			Name:    name,
			Props:   props,
			Content: content,
		}

		_, err = file.Insert(eng)
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}
	}

	response = fmt.Sprintf(tmplResponse, "ok", 0)

	return
}
