package accountancy

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"strings"
	"testing"
	"text/template"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tidwall/buntdb"
	"github.com/tidwall/gjson"
)

var dbDriver = "postgres"
var dbConn = "host=localhost port=5432 user=accountancy password=accountancy dbname=accountancy sslmode=disable"

func TestSystemUploadPostgres(t *testing.T) {

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

		response, err := Run(buff.String(), nil, nil)
		if err != nil {
			t.Error(err)
		}

		t.Error(response)

	}
}

func TestSystemUploadSqlite(t *testing.T) {

	names := []string{
		"/trait/terminal.js",
		"/trait/counterparty.js",
		"/trait/account.js",
		"/trait/transaction.js",
		"/trait/model.js",
		"/trait/client.js",
	}

	tmplFile, err := ioutil.ReadFile("./test/sqlite/upload.tmpl")
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

		err = ioutil.WriteFile("./test/sqlite/out/"+strings.Replace(name, "/", "-", -1)+".json", buff.Bytes(), 0755)
		if err != nil {
			t.Error(err)
		}

		response, err := Run(buff.String(), nil, nil)
		if err != nil {
			t.Error(err)
		}

		t.Error(response)

	}
}

func TestBatchSystemDropdbInitdbPostgres(t *testing.T) {

	src, err := ioutil.ReadFile("./test/dropdb-initdb.json")
	if err != nil {
		t.Error(err)
	}

	response, err := RunBatch(string(src), nil, nil)
	if err != nil {
		t.Error(err)
	}
	if response != `{"message": "ok", "status": 0}` {
		t.Error(response)
	}
}

func TestBatchSystemDropdbInitdbSqlite(t *testing.T) {

	src, err := ioutil.ReadFile("./test/sqlite/dropdb-initdb.json")
	if err != nil {
		t.Error(err)
	}

	response, err := RunBatch(string(src), nil, nil)
	if err != nil {
		t.Error(err)
	}

	t.Error(response)

}

func TestSystemInitdbPostgres(t *testing.T) {

	src, err := ioutil.ReadFile("./test/initdb.json")
	if err != nil {
		t.Error(err)
	}

	response, err := Run(string(src), nil, nil)
	if err != nil {
		t.Error(err)
	}
	if response != `{"message": "ok", "status": 0}` {
		t.Error(response)
	}
}

func TestSystemInitdbSqlite(t *testing.T) {

	src, err := ioutil.ReadFile("./test/sqlite/initdb.json")
	if err != nil {
		t.Error(err)
	}

	response, err := Run(string(src), nil, nil)
	if err != nil {
		t.Error(err)
	}
	if response != `{"message": "ok", "status": 0}` {
		t.Error(response)
	}
}

func TestSystemDropdbPostgres(t *testing.T) {

	src, err := ioutil.ReadFile("./test/dropdb.json")
	if err != nil {
		t.Error(err)
	}

	response, err := Run(string(src), nil, nil)
	if err != nil {
		t.Error(err)
	}
	if response != `{"message": "ok", "status": 0}` {
		t.Error(response)
	}
}

func TestSystemDropdbSqlite(t *testing.T) {

	src, err := ioutil.ReadFile("./test/sqlite/dropdb.json")
	if err != nil {
		t.Error(err)
	}

	response, err := Run(string(src), nil, nil)
	if err != nil {
		t.Error(err)
	}
	if response != `{"message": "ok", "status": 0}` {
		t.Error(response)
	}
}

func TestImportMetaPostgres(t *testing.T) {

	src, err := ioutil.ReadFile("./test/meta.json")
	if err != nil {
		t.Error(err)
	}

	response, err := Run(string(src), nil, nil)
	if err != nil {
		t.Error(err)
	}

	t.Error(response)

}

func TestImportMetaSqlite(t *testing.T) {

	src, err := ioutil.ReadFile("./test/sqlite/import-meta.json")
	if err != nil {
		t.Error(err)
	}

	response, err := Run(string(src), nil, nil)
	if err != nil {
		t.Error(err)
	}

	t.Error(response)

}

func TestImportDataSqlite(t *testing.T) {

	tmplFile, err := ioutil.ReadFile("./test/sqlite/import-data.tmpl")
	if err != nil {
		t.Error(err)
	}

	tmpl, err := template.New("data").Parse(string(tmplFile))
	if err != nil {
		t.Error(err)
	}

	traitNames := []interface{}{"Transaction"}

	db, err := buntdb.Open("../../../../../Niko/db/context-prod.db")
	if err != nil {
		t.Error(err)
		return
	}
	defer db.Close()

	eng, err := Connect(dbDriver, dbConn, true)
	if err != nil {
		t.Error(err)
		return
	}
	defer eng.Close()

	i := 0

	data := []map[string]interface{}{}

	err = db.View(func(tx *buntdb.Tx) (err error) {
		tx.Ascend("", func(key, value string) bool {
			i++
			//log.Println(i)

			if !gjson.Valid(value) {
				t.Errorf("%d - json is not valid, key: %s, value: %s", i, key, value)
				return true
			}

			dataJSON := gjson.Parse(value)
			name := dataJSON.Get("check.request.transactionId").String()
			props := dataJSON.Raw

			obj := map[string]interface{}{
				"Name":       name,
				"Props":      props,
				"TraitNames": traitNames,
			}

			data = append(data, obj)

			//if i > 50 {
			//	return false
			//}

			return true
		})
		return
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Error(i)

	var buff bytes.Buffer

	err = tmpl.Execute(&buff, data)
	if err != nil {
		t.Error(err)
	}

	err = ioutil.WriteFile("./test/sqlite/out/data.json", buff.Bytes(), 0755)
	if err != nil {
		t.Error(err)
	}

	//t.Error(buff.String())

	//return

	response, err := Run(buff.String(), nil, nil)
	if err != nil {
		t.Error(err)
	}

	t.Error(response)

}

func TestExportMetaSqlite(t *testing.T) {

	src, err := ioutil.ReadFile("./test/sqlite/export-meta.json")
	if err != nil {
		t.Error(err)
	}

	response, err := Run(string(src), nil, nil)
	if err != nil {
		t.Error(err)
	}

	t.Error(response)

}

func TestSelectObject(t *testing.T) {

	src, err := ioutil.ReadFile("./test/sqlite/select-object.json")
	if err != nil {
		t.Error(err)
	}

	response, err := Run(string(src), nil, nil)
	if err != nil {
		t.Error(err)
	}

	t.Error(response)

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
