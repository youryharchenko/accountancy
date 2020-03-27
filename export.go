package accountancy

import (
	"encoding/json"
	"fmt"

	"xorm.io/xorm"
)

// RunExport -
func RunExport(service string, request string) (response string, err error) {

	var eng *xorm.Engine
	eng, err = ConnectRequestDB(request)
	if err != nil {
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}
	defer eng.Close()

	switch service {
	case "meta":
		sess := eng.NewSession()
		defer sess.Close()

		err = sess.Begin()
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}

		response, err = ExportMeta(sess, request)
		if err != nil {
			err = sess.Rollback()
			if err != nil {
				response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			}
			return
		}

		err = sess.Commit()
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		}
	case "data":
		sess := eng.NewSession()
		defer sess.Close()

		err = sess.Begin()
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}

		response, err = ExportData(sess, request)
		if err != nil {
			err = sess.Rollback()
			if err != nil {
				response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			}
			return
		}

		err = sess.Commit()
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		}

	default:
		err = fmt.Errorf("RunExport: unknown service '%s'", service)
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
	}

	return
}

// ExportMeta -
func ExportMeta(sess *xorm.Session, request string) (response string, err error) {
	meta, err := LoadMeta(sess)
	if err != nil {
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}

	traits := []Trait{}
	for _, v := range meta.TraitsByID {
		traits = append(traits, *v)
	}
	relations := []Relation{}
	for _, v := range meta.RelationsByID {
		relations = append(relations, *v)
	}
	relationTraits := []RelationTraits{}
	for _, v := range meta.RelTraitsByID {
		relationTraits = append(relationTraits, *v)
	}
	mapBody := map[string]interface{}{
		"traits":          traits,
		"relations":       relations,
		"relation-traits": relationTraits,
	}

	mapResponse := map[string]interface{}{}
	err = json.Unmarshal([]byte(request), &mapResponse)
	if err != nil {
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}

	mapResponse["body"] = mapBody

	bufResponse, err := json.MarshalIndent(mapResponse, "", "  ")
	if err != nil {
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
		return
	}
	response = string(bufResponse)

	return
}

// ExportData -
func ExportData(sess *xorm.Session, request string) (response string, err error) {
	return
}
