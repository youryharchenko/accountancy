package accountancy

import (
	"fmt"
	"log"

	"github.com/tidwall/gjson"
	"xorm.io/core"
	"xorm.io/xorm"
)

// DB -
type DB interface {
	Insert(beans ...interface{}) (affected int64, err error)
	Get(bean interface{}) (has bool, err error)
	Find(rowsSlicePtr interface{}, condiBean ...interface{}) error
	Update(bean interface{}, condiBean ...interface{}) (affected int64, err error)
	ID(id interface{}) *xorm.Session
}

// Meta -
type Meta struct {
	TraitsByName    map[string]*Trait
	TraitsByID      map[int64]*Trait
	RelationsByName map[string]*Relation
	RelationsByID   map[int64]*Relation
	RelTraitsByName map[RelationTraitsKey]*RelationTraits
	RelTraitsByID   map[int64]*RelationTraits
	TypeDocsByName  map[string]*TypeDocument
	TypeDocsByID    map[int64]*TypeDocument
}

// LoadMeta -
func LoadMeta(db DB) (meta *Meta, err error) {
	meta = &Meta{
		TraitsByName:    map[string]*Trait{},
		TraitsByID:      map[int64]*Trait{},
		RelationsByName: map[string]*Relation{},
		RelationsByID:   map[int64]*Relation{},
		RelTraitsByName: map[RelationTraitsKey]*RelationTraits{},
		RelTraitsByID:   map[int64]*RelationTraits{},
	}
	err = meta.LoadRelations(db)
	if err != nil {
		return
	}
	err = meta.LoadTraits(db)
	if err != nil {
		return
	}
	err = meta.LoadRelTraits(db)
	if err != nil {
		return
	}

	log.Println("LoadMeta: loaded")
	return
}

// AddTrait -
func (meta *Meta) AddTrait(trait *Trait) {
	tmp := *trait
	meta.TraitsByName[trait.Name] = &tmp
	meta.TraitsByID[trait.ID] = &tmp
}

// AddRelation -
func (meta *Meta) AddRelation(relation *Relation) {
	tmp := *relation
	meta.RelationsByName[relation.Name] = &tmp
	meta.RelationsByID[relation.ID] = &tmp
}

// AddRelTraits -
func (meta *Meta) AddRelTraits(relTraits *RelationTraits) (err error) {
	tmp := *relTraits

	rel, ok := meta.RelationsByID[relTraits.RelationID]
	if !ok {
		err = fmt.Errorf("Meta.AddRelTraits: relation '%d' is missing", relTraits.RelationID)
		return
	}

	traitFrom, ok := meta.TraitsByID[relTraits.TraitFromID]
	if !ok {
		err = fmt.Errorf("Meta.AddRelTraits: trait from '%d' is missing", relTraits.TraitFromID)
		return
	}

	traitTo, ok := meta.TraitsByID[relTraits.TraitToID]
	if !ok {
		err = fmt.Errorf("Meta.AddRelTraits: trait from '%d' is missing", relTraits.TraitToID)
		return
	}

	tmp.RelationName = rel.Name
	tmp.TraitFromName = traitFrom.Name
	tmp.TraitToName = traitTo.Name

	meta.RelTraitsByID[relTraits.ID] = &tmp
	meta.RelTraitsByName[RelationTraitsKey{RelationName: rel.Name, TraitFromName: traitFrom.Name, TraitToName: traitTo.Name}] = &tmp

	return
}

// LoadTraits -
func (meta *Meta) LoadTraits(db DB) (err error) {
	list := []Trait{}
	err = db.Find(&list)
	if err != nil {
		return
	}

	for _, e := range list {
		jsFind := &File{Name: e.Lib}
		_, err = jsFind.Get(db)
		if err != nil {
			return
		}
		e.LibContent = string(jsFind.Content)
		meta.AddTrait(&e)
	}
	return
}

// LoadRelations -
func (meta *Meta) LoadRelations(db DB) (err error) {
	list := []Relation{}
	err = db.Find(&list)
	if err != nil {
		return
	}
	for _, e := range list {
		meta.AddRelation(&e)
	}
	return
}

// LoadRelTraits -
func (meta *Meta) LoadRelTraits(db DB) (err error) {
	list := []RelationTraits{}
	err = db.Find(&list)
	if err != nil {
		return
	}
	for _, e := range list {
		err = meta.AddRelTraits(&e)
		if err != nil {
			return
		}
	}
	return
}

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
		new(LinkObjects),
		new(Service),
		new(File),
	}
	return
}
