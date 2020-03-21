package accountancy

import (
	"fmt"

	"github.com/tidwall/gjson"
)

var tmplResponse = `{"message": "%s", "status": %d}`

// Run -
func Run(request string) (response string, err error) {
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
		response, err = RunSystem(servJSON.String(), request)
	case "import":
		if !servJSON.Exists() {
			err = fmt.Errorf("Run: import - field 'request.service' is missing")
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}
		response, err = RunImport(servJSON.String(), request)
	default:
		err = fmt.Errorf("Run: unknown command '%s'", command)
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}

	return
}

// RunBatch -
func RunBatch(request string) (response string, err error) {
	if !gjson.Valid(request) {
		err = fmt.Errorf("RunBatch: request is not valid JSON")
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}

	srcJSON := gjson.Parse(request)

	batchJSON := srcJSON.Array()
	for _, req := range batchJSON {
		response, err = Run(req.Raw)
		if err != nil {
			return
		}
	}

	return
}
