package model

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	couchdb "github.com/leesper/couchdb-golang"
	"golang.org/x/crypto/bcrypt"
)

//User What a User is
type User struct {
	_id       string            `json:"_id"`
	_rev      string            `json:"_rev"`
	Type      string            `json:"type"`
	Username  string            `json:"username"`
	Password  string            `json:"password"`
	Email     string            `json:"email"`
	Date      string            `json:"date"`
	Cardboxes []CardboxProgress `json:"cardboxes"`
	Boxnum    float64           `json:"boxnum"`
	Cardnum   float64           `json:"cardnum"`
	PicID     string            `json:"picid"`

	couchdb.Document
}

func GetUserDoc(userMap map[string]interface{}) (user User) {

	progBoxesMap := ToMapArray(userMap["cardboxes"])
	progBoxes := make([]CardboxProgress, len(progBoxesMap))

	for i, j := range progBoxesMap {
		progBoxes[i] = GetCardboxProgress(j)
	}

	user = User{
		_id:       userMap["_id"].(string),
		_rev:      userMap["_rev"].(string),
		Type:      userMap["type"].(string),
		Username:  userMap["username"].(string),
		Password:  userMap["password"].(string),
		Email:     userMap["email"].(string),
		Date:      userMap["date"].(string),
		Cardboxes: progBoxes,
		Boxnum:    userMap["boxnum"].(float64),
		Cardnum:   userMap["cardnum"].(float64),
		PicID:     userMap["picid"].(string),
	}
	return user
}

type CardboxProgress struct {
	BoxID    string   `json:"boxID"`
	Author   string   `json:"author"`
	Count    float64  `json:"count"`
	Lvl0     []string `json:"lvl0"`
	Lvl1     []string `json:"lvl1"`
	Lvl2     []string `json:"lvl2"`
	Lvl3     []string `json:"lvl3"`
	Lvl4     []string `json:"lvl4"`
	Done     []string `json:"done"`
	Progress float64  `json:"progress"`
}

func CreateNewCardboxProgress(boxID string) CardboxProgress {
	lvl0 := make([]string, 0)
	count := float64(0)
	cardbox, _ := GetCardboxByID(boxID)
	cards := ToMapArray(cardbox[0]["cards"])

	for i := 0; i < len(cards); i++ {
		lvl0 = append(lvl0, cards[i]["cardID"].(string))
		count++
	}

	newCardboxProgress := CardboxProgress{
		BoxID:    boxID,
		Count:    count,
		Lvl0:     lvl0,
		Lvl1:     make([]string, 0),
		Lvl2:     make([]string, 0),
		Lvl3:     make([]string, 0),
		Lvl4:     make([]string, 0),
		Done:     make([]string, 0),
		Author:   cardbox[0]["author"].(string),
		Progress: 0,
	}

	return newCardboxProgress
}
func GetCardboxProgress(mapProg map[string]interface{}) CardboxProgress {

	newCardboxProgress := CardboxProgress{
		BoxID:    mapProg["boxID"].(string),
		Count:    mapProg["count"].(float64),
		Lvl0:     makestringArray(mapProg["lvl0"].([]interface{})),
		Lvl1:     makestringArray(mapProg["lvl1"].([]interface{})),
		Lvl2:     makestringArray(mapProg["lvl2"].([]interface{})),
		Lvl3:     makestringArray(mapProg["lvl3"].([]interface{})),
		Lvl4:     makestringArray(mapProg["lvl4"].([]interface{})),
		Done:     makestringArray(mapProg["done"].([]interface{})),
		Author:   mapProg["author"].(string),
		Progress: mapProg["progress"].(float64),
	}

	return newCardboxProgress
}

func makestringArray(old []interface{}) []string {
	result := make([]string, len(old))
	for i, j := range old {
		result[i] = j.(string)
	}
	return result
}

//GetUser Gets User
func GetUser(username string) (map[string]interface{}, error) {
	user, err := btDB.QueryJSON(`
		{
   "selector": {
      "username": {
         "$eq": "` + username + `"
      }
   }
}`)
	if err != nil || len(user) == 0 {
		return map[string]interface{}{}, err
	}

	return user[0], nil
}

func GetIDbyUsername(username string) (string, error) {
	userdoc, err := btDB.QueryJSON(`
		{
   "selector": {
      "username": {
         "$eq": "bcamili"
      }
   }
}`)
	if err != nil {
		return "null", err
	}

	doc := userdoc[0]
	var docID string
	docID = doc["_id"].(string)
	return docID, err
}

func GetAllUsers() ([]map[string]interface{}, int, error) {
	allUsers, err := btDB.QueryJSON(`
		{
		   "selector": {
		      "username": {
		         "$gt": null
		      }
		   }
		}`)
	if err != nil {
		return nil, 0, err
	} else {
		return allUsers, len(allUsers), nil
	}
}

func (user User) Add() (err error) {
	// Check wether username already exists
	userInDB, err := GetUserByUsername(user.Username)
	if err == nil && userInDB.Username == user.Username {
		return errors.New("username exists already")
	}

	// Hash password
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	b64HashedPwd := base64.StdEncoding.EncodeToString(hashedPwd)
	user.PicID = base64.StdEncoding.EncodeToString([]byte(user.Username))
	user.Password = b64HashedPwd
	user.Type = "user"
	date := time.Now().Format("02.01.2006")
	user.Date = date
	user.Boxnum = float64(0)
	user.Cardnum = float64(0)
	user.Cardboxes = make([]CardboxProgress, 0)
	// Convert Todo struct to map[string]interface as required by Save() method
	u, err := user2Map(user)

	// Delete _id and _rev from map, otherwise DB access will be denied (unauthorized)
	delete(u, "_id")
	delete(u, "_rev")

	// Add todo to DB
	_, _, err = btDB.Save(u, nil)

	if err != nil {
		fmt.Printf("[Add] error 1: %s", err)
	}
	newFile, err := os.Create("static/docs/" + user.PicID)
	if err != nil {
		fmt.Printf("[Add] error 2: %s", err)
	}
	//	os.Chdir("/static/img/")
	defer newFile.Close()
	stdPic, err := ioutil.ReadFile("static/img/Profile-Pic.png")
	if err != nil {
		fmt.Printf("[Add] error 3: %s", err)
	}
	//os.Chdir("/static/docs/")
	_, err = newFile.Write(stdPic)

	if err != nil {
		fmt.Printf("[Add] error 4: %s", err)
	}
	//os.Chdir("/Users/benjamincamili/go/src/BrainTrain/")
	return err
}

// GetUserByUsername retrieve User by username
func GetUserByUsername(username string) (user User, err error) {
	if username == "" {
		return User{}, errors.New("no username provided")
	}
	userMaps, err := btDB.QueryJSON(`
		{
   "selector": {
      "username": {
         "$eq": "` + username + `"
      }
   }
}`)
	if err != nil || len(userMaps) != 1 {
		return User{}, err
	}

	user, err = map2User(userMaps[0])
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// Convert from User struct to map[string]interface{} as required by golang-couchdb methods
func user2Map(u User) (user map[string]interface{}, err error) {
	uJSON, err := json.Marshal(u)
	json.Unmarshal(uJSON, &user)

	return user, err
}

func map2User(user map[string]interface{}) (u User, err error) {
	uJSON, err := json.Marshal(user)
	json.Unmarshal(uJSON, &u)

	return u, err
}
