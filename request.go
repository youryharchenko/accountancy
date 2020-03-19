package accountancy

import (
	"fmt"

	"github.com/tidwall/gjson"
)

var tmplResponse = `{"message": "%s", "status": %d}`

// Run -
func Run(request string) (response string, err error) {
	if !gjson.Valid(request) {
		return "", fmt.Errorf("Run: request is not valid JSON")
	}

	srcJSON := gjson.Parse(request)

	reqJSON := srcJSON.Get("request")
	if !reqJSON.Exists() {
		return "", fmt.Errorf("Run: object 'request' is missing")
	}

	cmdJSON := reqJSON.Get("command")
	if !cmdJSON.Exists() {
		return "", fmt.Errorf("Run: field 'request.command' is missing")
	}

	command := cmdJSON.String()
	switch command {
	case "import":
		servJSON := reqJSON.Get("service")
		if !servJSON.Exists() {
			return "", fmt.Errorf("Run: import - field 'request.service' is missing")
		}
		response, err = RunImport(servJSON.String(), request)
	default:
		return "", fmt.Errorf("Run: unknown command '%s'", command)
	}

	return
}
