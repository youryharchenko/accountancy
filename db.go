package accountancy

import (
	"fmt"

	"github.com/tidwall/gjson"
	"xorm.io/core"
	"xorm.io/xorm"
)

// Connect -
func Connect(driverName string, connString string, show bool) (eng *xorm.Engine, err error) {
	eng, err = xorm.NewEngine(driverName, connString)
	if err != nil {
		return
	}
	eng.ShowSQL(show)
	eng.SetMapper(core.GonicMapper{})
	return
}

// ConnectRequestDB -
func ConnectRequestDB(request string) (eng *xorm.Engine, err error) {
	srcJSON := gjson.Parse(request)

	dbJSON := srcJSON.Get("db")
	if !dbJSON.Exists() {
		return nil, fmt.Errorf("ConnectRequestDB: object 'db' is missing")
	}

	drvJSON := dbJSON.Get("driver")
	if !drvJSON.Exists() {
		return nil, fmt.Errorf("ConnectRequestDB: field 'db.driver' is missing")
	}
	driver := drvJSON.String()

	connJSON := dbJSON.Get("connection")
	if !connJSON.Exists() {
		return nil, fmt.Errorf("ConnectRequestDB: field 'db.connection' is missing")
	}
	connection := connJSON.String()

	show := false
	showJSON := dbJSON.Get("show")
	if showJSON.Exists() {
		show = showJSON.Bool()
	}

	eng, err = Connect(driver, connection, show)

	return
}

// SyncAll -
func SyncAll(eng *xorm.Engine) (err error) {
	ents := entityList()

	for _, ent := range ents {
		err = Sync(eng, ent)
		if err != nil {
			return err
		}
	}

	return
}

// DropAll -
func DropAll(eng *xorm.Engine) (err error) {

	ents := entityList()
	err = eng.DropTables(ents...)

	return
}

// Sync -
func Sync(eng *xorm.Engine, ent interface{}) error {
	return eng.Sync2(ent)
}

func entityList() (ents []interface{}) {
	ents = []interface{}{
		new(Ledger),
		new(Object),
		new(Trait),
		new(TraitObject),
		new(Document),
		new(TypeDocument),
		new(Relation),
		new(RelationTraits),
		new(Service),
		new(File),
	}
	return
}
