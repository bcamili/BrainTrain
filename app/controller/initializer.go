package controller

import (
	"BrainTrain/app/model"
	"html/template"
)

var tpl *template.Template

type pageData struct {
	Title         string
	IsLoggedIn    bool
	MainInfo      MainInfo
	BoxName       string
	CatName       string
	SubcatName    string
	Userinfo      model.User
	ActiveCardbox model.Cardbox
	ActiveCards   []model.Card
	ActiveCard    model.Card
	Categories    []Category
	Levels        LevelDisplay
	ErrorMsg      Error
	Edit          bool
}

type Error struct {
	Error01 bool
	Error02 bool
	Error03 bool
	Error04 bool
}

type Category struct {
	CategoryName string
	Boxes        []model.Cardbox
}

func MapToCategory(oldMap map[string]interface{}) Category {
	boxes := model.ToMapArray(oldMap["Boxes"])
	boxesStruct := make([]model.Cardbox, len(boxes))
	for i := 0; i < len(boxes); i++ {
		boxesStruct[i], _ = model.GetCardbox(boxes[i])
	}
	newCategory := Category{
		CategoryName: oldMap["categoryname"].(string),
		Boxes:        boxesStruct,
	}

	return newCategory
}

type MainInfo struct {
	UserCount    int
	CardCount    float64
	BoxCount     int
	UserBoxCount int
}

type LevelDisplay struct {
	Lvl0            float64
	Lvl1            float64
	Lvl2            float64
	Lvl3            float64
	Lvl4            float64
	ActiveCardLevel string
	Progress        int
}

var stats MainInfo

func Init() {

	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))

}

func Stats() (stats MainInfo) {
	stats = MainInfo{}
	_, numUser, _ := model.GetAllUsers()
	Boxes, _ := model.GetPublicCardboxes()
	var numCards float64 = 0
	for i := 0; i < len(Boxes); i++ {
		numCards += Boxes[i]["numcards"].(float64)
	}
	stats.UserCount = numUser
	stats.CardCount = numCards
	stats.BoxCount = len(Boxes)
	return stats
}

func StatsU(username string) (stats MainInfo) {
	stats = Stats()
	user, _ := model.GetUser(username)
	userBoxes, _ := model.GetAllCardboxesByUsername(username)

	userBoxCount := 0

	for _, k := range user["cardboxes"].([]interface{}) {
		k2 := model.ToMap(k)
		if k2["author"].(string) != user["username"].(string) {
			userBoxCount++
		}
	}
	stats.UserBoxCount = userBoxCount + len(userBoxes)
	return stats
}
