package accountancy

import (
	"io/ioutil"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var dbDriver = "postgres"
var dbConn = "host=localhost port=5432 user=accountancy password=accountancy dbname=accountancy sslmode=disable"

func TestImportMeta(t *testing.T) {

	src, err := ioutil.ReadFile("./test/meta.json")
	if err != nil {
		t.Error(err)
	}

	response, err := Run(string(src))
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
