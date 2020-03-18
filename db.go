package lib

import (
	"github.com/go-xorm/xorm"
	"xorm.io/core"
)

// Connect -
func Connect(driverName string, connString string) (eng *xorm.Engine, err error) {
	eng, err = xorm.NewEngine(driverName, connString)
	//eng.SetLogLevel(core.LOG_WARNING)
	eng.ShowSQL(true)
	eng.SetMapper(core.GonicMapper{})
	return
}

// SyncAll -
func SyncAll(eng *xorm.Engine) (err error) {
	ents := []interface{}{
		new(Ledger),
		new(Object),
		new(Trait),
		new(TraitObject),
		new(Document),
		new(TypeDocument),
		new(Relation),
		new(RelationTraits),
		new(Service),
	}

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

	ents := []interface{}{
		new(Ledger),
		new(Object),
		new(Trait),
		new(TraitObject),
		new(Document),
		new(TypeDocument),
		new(Relation),
		new(RelationTraits),
		new(Service),
	}

	err = eng.DropTables(ents...)

	return
}

// Sync -
func Sync(eng *xorm.Engine, ent interface{}) error {
	return eng.Sync2(ent)
}
