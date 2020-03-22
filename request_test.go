package accountancy

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"testing"
	"text/template"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var dbDriver = "postgres"
var dbConn = "host=localhost port=5432 user=accountancy password=accountancy dbname=accountancy sslmode=disable"

func TestSystemUpload(t *testing.T) {

	names := []string{
		"/trait/terminal.js",
		"/trait/counterparty.js",
		"/trait/account.js",
		"/trait/transaction.js",
	}

	tmplFile, err := ioutil.ReadFile("./test/upload.tmpl")
	if err != nil {
		t.Error(err)
	}

	tmpl, err := template.New("upload").Parse(string(tmplFile))
	if err != nil {
		t.Error(err)
	}

	for _, name := range names {

		contFile, err := ioutil.ReadFile("./test/file" + name)
		if err != nil {
			t.Error(err)
			continue
		}

		content := base64.StdEncoding.EncodeToString(contFile)

		data := map[string]interface{}{
			"UUID":    uuid.New().String(),
			"Name":    name,
			"Content": content,
		}
		var buff bytes.Buffer

		err = tmpl.Execute(&buff, data)
		if err != nil {
			t.Error(err)
		}

		response, err := Run(buff.String())
		if err != nil {
			t.Error(err)
		}

		if response != `{"message": "ok", "status": 0}` {
			t.Error(response)
		}
	}
}

func TestBatchSystemDropdbInitdb(t *testing.T) {

	src, err := ioutil.ReadFile("./test/dropdb-initdb.json")
	if err != nil {
		t.Error(err)
	}

	response, err := RunBatch(string(src))
	if err != nil {
		t.Error(err)
	}
	if response != `{"message": "ok", "status": 0}` {
		t.Error(response)
	}
}

func TestSystemInitdb(t *testing.T) {

	src, err := ioutil.ReadFile("./test/initdb.json")
	if err != nil {
		t.Error(err)
	}

	response, err := Run(string(src))
	if err != nil {
		t.Error(err)
	}
	if response != `{"message": "ok", "status": 0}` {
		t.Error(response)
	}
}

func TestSystemDropdb(t *testing.T) {

	src, err := ioutil.ReadFile("./test/dropdb.json")
	if err != nil {
		t.Error(err)
	}

	response, err := Run(string(src))
	if err != nil {
		t.Error(err)
	}
	if response != `{"message": "ok", "status": 0}` {
		t.Error(response)
	}
}

func TestImportMeta1(t *testing.T) {

	src, err := ioutil.ReadFile("./test/meta.json")
	if err != nil {
		t.Error(err)
	}

	response, err := Run(string(src))
	if err != nil {
		t.Error(err)
	}
	if response != `{"message": "ok", "status": 0}` {
		t.Error(response)
	}
}

func TestSync(t *testing.T) {

	eng, err := Connect(dbDriver, dbConn, true)
	if err != nil {
		t.Error(err)
	}
	defer eng.Close()

	err = SyncAll(eng)
	if err != nil {
		t.Error(err)
	}

}

func TestDrop1(t *testing.T) {

	eng, err := Connect(dbDriver, dbConn, true)
	if err != nil {
		t.Error(err)
	}
	defer eng.Close()

	err = DropAll(eng)
	if err != nil {
		t.Error(err)
	}

}

func TestUUID(t *testing.T) {
	t.Error(uuid.New().String())
}
