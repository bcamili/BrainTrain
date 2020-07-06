package model

import (
	"encoding/json"
	"fmt"
)

//Cardbox What a Cardbox is
type Cardbox struct {
	Id          string  `json:"_id"`
	Type        string  `json:"type"`
	Boxname     string  `json:"boxname"`
	NumCards    float64 `json:"numcards"`
	Author      string  `json:"author"`
	Category    string  `json:"category"`
	Subcategory string  `json:"subcategory"`
	Description string  `json:"description"`
	Cards       []Card  `json:"cards"`
	Visibility  string  `json:"visibility"`
	Select      string
	Progress    float64
}

//GetCardbox gets Cardbox
func GetCardbox(t map[string]interface{}) (Cardbox, error) {

	cardsFace := t["cards"].([]interface{})
	cards := make([]Card, len(cardsFace))

	for i, j := range cardsFace {
		cards[i], _ = GetCard(ToMap(j))
	}

	todo := Cardbox{
		Id:          t["_id"].(string),
		Type:        "cardbox",
		Boxname:     t["boxname"].(string),
		NumCards:    t["numcards"].(float64),
		Author:      t["author"].(string),
		Category:    t["category"].(string),
		Subcategory: t["subcategory"].(string),
		Description: t["description"].(string),
		Cards:       cards,
		Visibility:  t["visibility"].(string),
	}
	return todo, nil
}

func (box Cardbox) Add() (id string, err error) {
	box.Type = "cardbox"
	box.NumCards = float64(0)
	box.Cards = make([]Card, 0)
	// Convert Todo struct to map[string]interface as required by Save() method
	b, err := Cardbox2Map(box)

	// Delete _id and _rev from map, otherwise DB access will be denied (unauthorized)
	delete(b, "_id")
	delete(b, "_rev")
	delete(b, "Progress")
	delete(b, "Select")
	// Add todo to DB
	id, _, err = btDB.Save(b, nil)

	if err != nil {
		fmt.Printf("[Add] error: %s", err)
	}

	return id, err
}

func GetAllCardboxes() ([]map[string]interface{}, error) {
	boxes, err := btDB.QueryJSON(`
			{
	   "selector": {
	      "type": {
	         "$eq": "box"
	      }
	   }
	}`)
	if err != nil {
		return nil, err
	}

	return boxes, nil
}

func GetAllCardboxesByUsername(username string) ([]map[string]interface{}, error) {
	boxes, err := btDB.QueryJSON(`
			{
	   "selector": {
	      "author": {
	         "$eq": "` + username + `"
	      }
	   }
	}`)
	if err != nil {
		return nil, err
	}

	return boxes, nil
}

func GetCardboxByID(ID string) ([]map[string]interface{}, error) {
	boxes, err := btDB.QueryJSON(`
			{
	   "selector": {
	      "_id": {
	         "$eq": "` + ID + `"
	      }
	   }
	}`)
	if err != nil {
		return nil, err
	}

	return boxes, nil
}

func GetPublicCardboxes() ([]map[string]interface{}, error) {
	boxes, err := btDB.QueryJSON(
		`{
	   "selector": {
	      "visibility": {
	         "$eq": "Public"
	      }
	   }
	}`)

	if err != nil {
		return nil, err
	}
	return boxes, nil
}

func GetSubcatCardboxes(cat string, subcat string) ([]Cardbox, error) {
	boxes, err := btDB.QueryJSON(
		`{
			"selector": {
       "type": {
          "$eq": "cardbox"
       },
			 "visibility": {
          "$eq": "Public"
       },
       "category": {
          "$eq": "` + cat + `"
       },
			 "subcategory": {
          "$eq": "` + subcat + `"
       }
    }
	}`)

	if err != nil {
		return nil, err
	}

	cardboxes := make([]Cardbox, 0)
	for _, v := range boxes {
		box, _ := GetCardbox(v)
		cardboxes = append(cardboxes, box)
	}
	return cardboxes, nil
}

func GetUserCategories(username string) ([]map[string]interface{}, error) {
	user, err := GetUser(username)
	if err != nil {
		var empty []map[string]interface{}
		return empty, err
	}
	categories := ToMapArray(user["usercardboxes"])
	return categories, nil
}

func GetBoxfromCategories(cats map[string]interface{}) []map[string]interface{} {
	subcat := ToMapArray(cats["subcategories"])
	var boxes []map[string]interface{}
	for i := 0; i < len(subcat); i++ {
		cbs := ToMapArray(subcat[i]["cardboxes"])
		for j := 0; j < len(cbs); j++ {
			boxes = append(boxes, cbs[j])
		}

	}
	return boxes
}

func GetBoxesWithOwner(owner string) []map[string]interface{} {
	user, _ := GetUser(owner)

	var userCBs []map[string]interface{}

	a := ToMapArray(user["usercardboxes"])

	for i := 0; i < len(a); i++ {
		b := ToMapArray(a[i]["subcategories"])
		for j := 0; j < len(b); j++ {
			c := ToMapArray(b[j]["cardboxes"])
			for k := 0; k < len(c); k++ {
				if c[k]["owner"].(string) == owner {
					userCBs = append(userCBs, c[k])
				}
			}
		}
	}
	return userCBs
}

//AddCardbox adds Cardbox
func AddCardbox(bx Cardbox, bxs []Cardbox) []Cardbox {
	bxs = append(bxs, bx)
	return bxs
}

func Cardbox2Map(c Cardbox) (box map[string]interface{}, err error) {
	uJSON, err := json.Marshal(c)
	json.Unmarshal(uJSON, &box)

	return box, err
}
