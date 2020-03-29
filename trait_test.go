package accountancy

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/tidwall/buntdb"
	"github.com/tidwall/gjson"
)

func TestCreateTrx(t *testing.T) {

	traitName := "Transaction"

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

	//keys := []string{}
	//values := []string{}
	i := 0

	err = db.View(func(tx *buntdb.Tx) (err error) {
		tx.Ascend("", func(key, value string) bool {
			i++
			log.Println(i)
			//log.Println(key, value)

			//keys = append(keys, key)
			//values = append(values, value)
			if !gjson.Valid(value) {
				t.Errorf("%d - json is not valid, key: %s, value: %s", i, key, value)
				return true
			}
			//ctxJSON := gjson.Parse(value)

			ctx := map[string]interface{}{}

			if err := json.Unmarshal([]byte(value), &ctx); err != nil {
				t.Error(i, err)
				return true
			}

			obj, err := NewObject(key, ctx)
			if err != nil {
				t.Error(i, err)
				return true
			}

			sess := eng.NewSession()
			defer sess.Close()

			traitTrx := &Trait{Name: traitName}
			ok, err := traitTrx.Get(sess)
			if err != nil {
				t.Error(err)
				err = sess.Rollback()
				if err != nil {
					t.Error(i, err)
				}
				return true
			}
			if !ok {
				t.Errorf("trait '%s' is missing", traitName)
				err = sess.Rollback()
				if err != nil {
					t.Error(i, err)
				}
				return true
			}

			_, _, err = obj.InsertOrUpdate(sess, &Object{Name: key})
			if err != nil {
				t.Error(i, err)
				err = sess.Rollback()
				if err != nil {
					t.Error(i, err)
				}
				return true
			}

			_, _, err = traitTrx.AddObject(sess, obj, nil)
			if err != nil {
				t.Error(i, err)
				err = sess.Rollback()
				if err != nil {
					t.Error(i, err)
				}
				return true
			}

			err = sess.Commit()
			if err != nil {
				t.Error(i, err)
			}
			//if i > 10 {
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
}
