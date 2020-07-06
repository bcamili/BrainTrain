package model

import "html/template"

//Card What a Card is
type Card struct {
	CardID       string `json:"cardID"`
	Title        string `json:"title"`
	Question     string `json:"question"`
	QuestionHTML template.HTML
	Answer       string `json:"answer"`
	AnswerHTML   template.HTML

	Select     string
	Levelmeter []LevelMeter
	Count      int
}

type LevelMeter struct {
	Level     string
	LevelCode string
}

//GetCard gets Card
func GetCard(t map[string]interface{}) (Card, error) {

	todo := Card{
		CardID:   t["cardID"].(string),
		Title:    t["title"].(string),
		Question: t["question"].(string),
		Answer:   t["answer"].(string),
	}
	return todo, nil
}

func GetAllCards() ([]map[string]interface{}, int) {
	allBoxes, _ := GetAllCardboxes()
	var allCards []map[string]interface{}
	for k := 0; k < len(allBoxes); k++ {
		cards := ToMapArray(allBoxes[k]["cards"])
		for l := 0; l < len(cards); l++ {
			card := cards[l]
			allCards = append(allCards, card)
		}
	}
	return allCards, len(allCards)
}

//AddCard adds Card
func AddCard(c Card, box MapArray) MapArray {
	// Convert Todo struct to map[string]interface as required by Save() method
	newCard := Map{
		"cardID":   c.CardID,
		"title":    c.Title,
		"question": c.Question,
		"answer":   c.Answer,
	}
	box = append(box, newCard)

	return box
}
