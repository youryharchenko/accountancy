package accountancy

import (
	"fmt"
	"log"

	"github.com/tidwall/gjson"
	"xorm.io/xorm"
)

var tmplResponse = `{"message": "%s", "status": %d}`

// Run -
func Run(request string, meta *Meta, db DB) (response string, err error) {
	srcJSON := gjson.Parse(request)
	dbJSON := srcJSON.Get("db")
	reqJSON := srcJSON.Get("request")
	bodyJSON := srcJSON.Get("body")
	return RunBatch(fmt.Sprintf(`{"db":%s,"batch":[{"request":%s,"body":%s}]}`, dbJSON.Raw, reqJSON.Raw, bodyJSON.Raw), meta, db)
}

// RunOne -
func RunOne(request string, meta *Meta, db DB) (response string, err error) {

	if !gjson.Valid(request) {
		err = fmt.Errorf("Run: request is not valid JSON")
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}

	srcJSON := gjson.Parse(request)

	reqJSON := srcJSON.Get("request")
	if !reqJSON.Exists() {
		err = fmt.Errorf("Run: object 'request' is missing")
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}

	cmdJSON := reqJSON.Get("command")
	if !cmdJSON.Exists() {
		err = fmt.Errorf("Run: field 'request.command' is missing")
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}

	servJSON := reqJSON.Get("service")

	command := cmdJSON.String()
	switch command {
	case "system":
		if !servJSON.Exists() {
			err = fmt.Errorf("Run: system - field 'request.service' is missing")
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}
		response, err = RunSystem(servJSON.String(), request, meta, db)
	case "import":
		if !servJSON.Exists() {
			err = fmt.Errorf("Run: import - field 'request.service' is missing")
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}
		response, err = RunImport(servJSON.String(), request, meta, db)
	case "export":
		if !servJSON.Exists() {
			err = fmt.Errorf("Run: export - field 'request.service' is missing")
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}
		response, err = RunExport(servJSON.String(), request, meta, db)
	case "insert":
		if !servJSON.Exists() {
			err = fmt.Errorf("Run: insert - field 'request.service' is missing")
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}
		response, err = RunInsert(servJSON.String(), request, meta, db)
	default:
		err = fmt.Errorf("Run: unknown command '%s'", command)
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}

	return
}

// RunBatch -
func RunBatch(request string, meta *Meta, db DB) (response string, err error) {
	log.Println("RunBatch:", request)

	if !gjson.Valid(request) {
		err = fmt.Errorf("RunBatch: request is not valid JSON")
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}

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
		defer func() {
			if !sess.IsClosed() {
				sess.Commit()
			}
		}()

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
			defer func() {
				if !sess.IsClosed() {
					sess.Commit()
				}
			}()
		case *xorm.Session:
			sess = db.(*xorm.Session)
		}
	}

	srcJSON := gjson.Parse(request).Get("batch")
	batchJSON := srcJSON.Array()
	for _, req := range batchJSON {
		response, err = RunOne(req.Raw, meta, sess)
		if err != nil {
			sess.Rollback()
			sess.Close()
			return
		}
	}

	return
}
