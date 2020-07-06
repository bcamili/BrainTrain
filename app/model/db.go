package model

import (
	"encoding/json"
	"fmt"

	couchdb "github.com/leesper/couchdb-golang"
)

type MapArray []map[string]interface{}
type Map map[string]interface{}

// Db handle
var btDB *couchdb.Database

//Init initializes DB
func Init() {
	var err error
	btDB, err = couchdb.NewDatabase("http://127.0.0.1:5984/todo")
	if err != nil {
		panic(err)
	}

}
func UpdateFile(f map[string]interface{}, fid string) (id string, err error) {
	_ = btDB.Delete(fid)
	delete(f, "_id")
	delete(f, "_rev")
	id, _, err = btDB.Save(f, nil)

	if err != nil {
		fmt.Printf("[UpdateFile] error: %s", err)
	}
	return id, err
}

// Delete Todo with the provided id from DB
func DeleteFile(id string) error {
	err := btDB.Delete(id)

	return err
}

// ---------------------------------------------------------------------------
// Internal helper functions
// ---------------------------------------------------------------------------

// Convert from Todo struct to map[string]interface{} as required by Set() method
// func todo2Map(t Todo) map[string]interface{} {
// 	var doc map[string]interface{}
// 	tJSON, _ := json.Marshal(t)
// 	json.Unmarshal(tJSON, &doc)
//
// 	return doc
// }

//GetAllTodos Get all Todos from DB
func GetAllTodos() ([]map[string]interface{}, error) {
	allTodos, err := btDB.QueryJSON(`
		{
		   "selector": {
		      "_id": {
		         "$gt": null
		      }
		   }
		}`)
	if err != nil {
		return nil, err
	} else {
		return allTodos, nil
	}
}
func ToMapArray(i interface{}) []map[string]interface{} {
	var doc []map[string]interface{}
	tj, _ := json.Marshal(i)
	json.Unmarshal(tj, &doc)
	return doc
}
func ToMap(i interface{}) map[string]interface{} {
	var doc map[string]interface{}
	tj, _ := json.Marshal(i)
	json.Unmarshal(tj, &doc)
	return doc
}
