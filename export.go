package accountancy

import (
	"encoding/json"
	"fmt"
)

// RunExport -
func RunExport(service string, request string, meta *Meta, db DB) (response string, err error) {

	switch service {
	case "meta":
		response, err = ExportMeta(db, request)
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}
	case "data":
		response, err = ExportData(db, request)
		if err != nil {
			response = fmt.Sprintf(tmplResponse, err.Error(), -1)
			return
		}
	default:
		err = fmt.Errorf("RunExport: unknown service '%s'", service)
		response = fmt.Sprintf(tmplResponse, err.Error(), -1)
	}

	return
}

// ExportMeta -
func ExportMeta(db DB, request string) (response string, err error) {
	meta, err := LoadMeta(db)
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
func ExportData(db DB, request string) (response string, err error) {
	return
}
