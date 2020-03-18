package main

import (
	"io/ioutil"
	"testing"

	l "github.com/youryharchenko/accountancy/lib"

	_ "github.com/lib/pq"
)

var dbDriver = "postgres"
var dbConn = "host=localhost port=5432 user=accountancy password=accountancy dbname=accountancy sslmode=disable"

func TestSync(t *testing.T) {

	eng, err := l.Connect(dbDriver, dbConn)
	if err != nil {
		t.Error(err)
	}
	defer eng.Close()

	err = l.SyncAll(eng)
	if err != nil {
		t.Error(err)
	}

}

func TestDrop(t *testing.T) {

	eng, err := l.Connect(dbDriver, dbConn)
	if err != nil {
		t.Error(err)
	}
	defer eng.Close()

	err = l.DropAll(eng)
	if err != nil {
		t.Error(err)
	}

}

func TestNewTrait(t *testing.T) {
	eng, err := l.Connect(dbDriver, dbConn)
	if err != nil {
		t.Error(err)
	}
	defer eng.Close()

	trait, err := l.NewTrait("Terminal", map[string]interface{}{"address": ""}, "")
	if err != nil {
		t.Error(err)
	}

	_, err = trait.Insert(eng)
	if err != nil {
		t.Error(err)
	}

}
func TestUpdateTrait(t *testing.T) {
	eng, err := l.Connect(dbDriver, dbConn)
	if err != nil {
		t.Error(err)
	}
	defer eng.Close()

	trait := &l.Trait{ID: 3}
	has, err := trait.Get(eng)
	if err != nil {
		t.Error(err)
	}
	if !has {
		t.Error("Get 'has' == false")
	}

	trait.Name = "Counterparty"

	_, err = trait.Update(eng)
	if err != nil {
		t.Error(err)
	}

}

func TestTraitAddObject(t *testing.T) {
	eng, err := l.Connect(dbDriver, dbConn)
	if err != nil {
		t.Error(err)
	}
	defer eng.Close()

	object := &l.Object{ID: 2}
	has, err := object.Get(eng)
	if err != nil {
		t.Error(err)
	}
	if !has {
		t.Error("Get 'has' == false")
	}

	trait := &l.Trait{ID: 4}
	has, err = trait.Get(eng)
	if err != nil {
		t.Error(err)
	}
	if !has {
		t.Error("Get 'has' == false")
	}

	aff, err := trait.AddObject(eng, object)
	if err != nil {
		t.Error(err)
	}
	if aff != 1 {
		t.Error(aff)
	}

}

func TestTraitUpdateObjectProps(t *testing.T) {
	eng, err := l.Connect(dbDriver, dbConn)
	if err != nil {
		t.Error(err)
	}
	defer eng.Close()

	object := &l.Object{ID: 2}
	has, err := object.Get(eng)
	if err != nil {
		t.Error(err)
	}
	if !has {
		t.Error("Get 'has' == false")
	}

	trait := &l.Trait{ID: 4}
	has, err = trait.Get(eng)
	if err != nil {
		t.Error(err)
	}
	if !has {
		t.Error("Get 'has' == false")
	}

	aff, err := trait.UpdateObjectProps(eng, object, map[string]interface{}{"address": "м.Львів, вул.Лисенка, 14-б"})
	if err != nil {
		t.Error(err)
	}
	if aff != 1 {
		t.Error(aff)
	}

}

func TestNewObject(t *testing.T) {
	eng, err := l.Connect(dbDriver, dbConn)
	if err != nil {
		t.Error(err)
	}
	defer eng.Close()

	object, err := l.NewObject("1129", map[string]interface{}{})
	if err != nil {
		t.Error(err)
	}

	_, err = object.Insert(eng)
	if err != nil {
		t.Error(err)
	}

}

func TestUpdateObject(t *testing.T) {
	eng, err := l.Connect(dbDriver, dbConn)
	if err != nil {
		t.Error(err)
	}
	defer eng.Close()

	object := &l.Object{ID: 1}
	has, err := object.Get(eng)
	if err != nil {
		t.Error(err)
	}
	if !has {
		t.Error("Get 'has' == false")
	}

	object.Name = "1124"

	_, err = object.Update(eng)
	if err != nil {
		t.Error(err)
	}

}

func TestObjectAddTrait(t *testing.T) {
	eng, err := l.Connect(dbDriver, dbConn)
	if err != nil {
		t.Error(err)
	}
	defer eng.Close()

	object := &l.Object{ID: 2}
	has, err := object.Get(eng)
	if err != nil {
		t.Error(err)
	}
	if !has {
		t.Error("Get 'has' == false")
	}

	trait := &l.Trait{ID: 4}
	has, err = trait.Get(eng)
	if err != nil {
		t.Error(err)
	}
	if !has {
		t.Error("Get 'has' == false")
	}

	aff, err := object.AddTrait(eng, trait)
	if err != nil {
		t.Error(err)
	}
	if aff != 1 {
		t.Error(aff)
	}

}

func TestNewService(t *testing.T) {
	eng, err := l.Connect(dbDriver, dbConn)
	if err != nil {
		t.Error(err)
	}
	defer eng.Close()

	//serviceName := "megapolissto"
	//serviceName := "nparts"
	serviceName := "invoice"

	js, err := ioutil.ReadFile("../script/service/" + serviceName + ".js")
	if err != nil {
		t.Error(err)
	}

	service, err := l.NewService(serviceName, map[string]interface{}{}, string(js))
	if err != nil {
		t.Error(err)
	}

	_, err = service.Insert(eng)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateServiceLib(t *testing.T) {
	eng, err := l.Connect(dbDriver, dbConn)
	if err != nil {
		t.Error(err)
	}
	defer eng.Close()

	serviceName := "megapolissto"
	//serviceName := "nparts"
	//serviceName := "invoice"

	service := &l.Service{Name: serviceName}
	has, err := service.Get(eng)
	if err != nil {
		t.Error(err)
	}
	if !has {
		t.Error("Get 'has' == false")
	}

	js, err := ioutil.ReadFile("../script/service/" + serviceName + ".js")
	if err != nil {
		t.Error(err)
	}

	service.Lib = string(js)

	_, err = service.Update(eng)
	if err != nil {
		t.Error(err)
	}

}

func TestImportFromService(t *testing.T) {
	eng, err := l.Connect(dbDriver, dbConn)
	if err != nil {
		t.Error(err)
	}
	defer eng.Close()

	ctx, err := ioutil.ReadFile("../test/ctx-megapolis.json")
	if err != nil {
		t.Error(err)
	}

	_, err = l.ImportFromService(eng, string(ctx))
	if err != nil {
		t.Error(err)
	}
}
