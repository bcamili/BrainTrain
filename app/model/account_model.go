package model

import (
	"fmt"

	couchdb "github.com/leesper/couchdb-golang"
)

//Account What an account is
type Account struct {
	Type     string `json:"type"`
	PID      string `json:"pid"`
	Email    string `json:"email"`
	Password string `json:"password"`
	couchdb.Document
}

//GetAccount Gets account
func GetAccount(id string) (Account, error) {
	t, err := btDB.Get(id, nil)
	if err != nil {
		return Account{}, err
	}
	todo := Account{
		Type:     t["type"].(string),
		PID:      t["pid"].(string),
		Email:    t["email"].(string),
		Password: t["password"].(string),
	}
	return todo, nil
}

//AddAccount Adds account
func (a Account) AddAccount() error {
	// Convert Todo struct to map[string]interface as required by Save() method
	todo, _ := couchdb.ToJSONCompatibleMap(a)
	// Delete _id and _rev from map, otherwise DB access will be denied (unauthorized)
	delete(todo, "_id")
	delete(todo, "_rev")
	// Add todo to DB
	_, _, err := btDB.Save(todo, nil)

	if err != nil {
		fmt.Printf("[AddAccount] error: %s", err)
	}
	return err
}
