package controller

import (
	"BrainTrain/app/model"
	"html/template"
	"log"
	"math/rand"
	"net/http"

	"github.com/russross/blackfriday"
)

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func Lrn(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "session")
	pd := session.Values["pagedata"].(pageData)

	user, _ := model.GetUser(pd.Userinfo.Username)
	userboxes := model.ToMapArray(user["cardboxes"])
	var progressBox map[string]interface{}

	isWorkedOn := false
	for i := 0; i < len(userboxes); i++ {
		if userboxes[i]["boxID"] == pd.ActiveCardbox.Id {

			isWorkedOn = true
			progressBox = userboxes[i]
		}
	}

	if !isWorkedOn {
		newBox := model.CreateNewCardboxProgress(pd.ActiveCardbox.Id)
		newBoxMap := model.ToMap(newBox)
		userboxes = append(userboxes, newBoxMap)
		user["cardboxes"] = userboxes
		model.UpdateFile(user, user["_id"].(string))
		progressBox = newBoxMap

	}

	randomStack := getRandomStackNumber()
	idStack := progressBox[randomStack].([]interface{})
	isThereEvenAny := 0
	for len(idStack) == 0 {
		randomStack = getRandomStackNumber()
		idStack = progressBox[randomStack].([]interface{})
		isThereEvenAny++
		if isThereEvenAny > 5 {
			pd.Title = "BrainTrain: Lernen"
			pd.MainInfo = StatsU(pd.Userinfo.Username)
			session.Values["pagedata"] = pd
			session.Save(req, w)
			err := tpl.ExecuteTemplate(w, "learn.gohtml", pd)

			if err != nil {
				log.Println(err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			return
		}
	}
	cardID := idStack[rand.Intn(len(idStack))].(string)

	for i := 0; i < len(pd.ActiveCards); i++ {
		if pd.ActiveCards[i].CardID == cardID {
			pd.ActiveCard = pd.ActiveCards[i]
			pd.Levels.ActiveCardLevel = randomStack
		}
	}

	pd.Levels.Lvl0 = float64(len(progressBox["lvl0"].([]interface{})))
	pd.Levels.Lvl1 = float64(len(progressBox["lvl1"].([]interface{})))
	pd.Levels.Lvl2 = float64(len(progressBox["lvl2"].([]interface{})))
	pd.Levels.Lvl3 = float64(len(progressBox["lvl3"].([]interface{})))
	pd.Levels.Lvl4 = float64(len(progressBox["lvl4"].([]interface{})))

	allCards := pd.Levels.Lvl0 + pd.Levels.Lvl1 + pd.Levels.Lvl2 + pd.Levels.Lvl3 + pd.Levels.Lvl4

	pd.Levels.Progress = int((pd.Levels.Lvl1 + pd.Levels.Lvl2*float64(2) + pd.Levels.Lvl3*float64(3) + pd.Levels.Lvl4*float64(4)) * float64(25) / allCards)

	for i := 0; i < len(userboxes); i++ {
		if userboxes[i]["boxID"] == pd.ActiveCardbox.Id {
			progressBox = userboxes[i]
			progressBox["progress"] = pd.Levels.Progress
			userboxes[i] = progressBox

		}
	}

	user["cardboxes"] = userboxes
	model.UpdateFile(user, user["_id"].(string))

	userMap, _ := model.GetUser(pd.Userinfo.Username)
	user2 := model.GetUserDoc(userMap)
	userCardboxes := user2.Cardboxes
	for _, i := range userCardboxes {
		if i.BoxID == pd.ActiveCardbox.Id {
			pd.ActiveCardbox.Progress = i.Progress

			lvls := make([]model.LevelMeter, 5)
			lvls[0].Level = "0"
			lvls[1].Level = "1"
			lvls[2].Level = "2"
			lvls[3].Level = "3"
			lvls[4].Level = "4"
			for i := 0; i < 5; i++ {
				lvls[i].LevelCode = "none"
			}

			if contains(i.Lvl0, pd.ActiveCard.CardID) {
				lvls[0].LevelCode = "currentLevel"
			}
			if contains(i.Lvl1, pd.ActiveCard.CardID) {
				lvls[1].LevelCode = "currentLevel"
			}
			if contains(i.Lvl2, pd.ActiveCard.CardID) {
				lvls[2].LevelCode = "currentLevel"
			}
			if contains(i.Lvl3, pd.ActiveCard.CardID) {
				lvls[3].LevelCode = "currentLevel"
			}
			if contains(i.Lvl4, pd.ActiveCard.CardID) {
				lvls[4].LevelCode = "currentLevel"
			}

			pd.ActiveCard.Levelmeter = lvls
		}
	}
	pd.ActiveCard.QuestionHTML = template.HTML(blackfriday.MarkdownCommon([]byte(pd.ActiveCard.Question)))
	pd.ActiveCard.AnswerHTML = template.HTML(blackfriday.MarkdownCommon([]byte(pd.ActiveCard.Answer)))

	pd.Title = "BrainTrain: Lernen"
	pd.MainInfo = StatsU(pd.Userinfo.Username)
	session.Values["pagedata"] = pd
	session.Save(req, w)

	err := tpl.ExecuteTemplate(w, "learn.gohtml", pd)

	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func getRandomStackNumber() string {
	randomStack := rand.Intn(14)
	randomStackNumber := "lvl0"

	if randomStack == 0 {
		randomStackNumber = "lvl4"
	}
	if randomStack > 0 && randomStack < 3 {
		randomStackNumber = "lvl3"
	}
	if randomStack >= 3 && randomStack < 6 {
		randomStackNumber = "lvl2"
	}
	if randomStack >= 6 && randomStack < 10 {
		randomStackNumber = "lvl1"
	}

	return randomStackNumber
}

func Lrn2(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "session")
	pd := session.Values["pagedata"].(pageData)
	if session.Values["authenticated"] != nil {
		pd.IsLoggedIn = session.Values["authenticated"].(bool)
	} else {
		pd.IsLoggedIn = false
	}
	pd.Title = "BrainTrain: Lernen"
	pd.MainInfo = StatsU(pd.Userinfo.Username)
	session.Values["pagedata"] = pd
	session.Save(req, w)
	err := tpl.ExecuteTemplate(w, "learn-2.gohtml", pd)

	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func LearnBox(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "session")
	pd := session.Values["pagedata"].(pageData)
	boxID := req.URL.EscapedPath()[11:len(req.URL.EscapedPath())]

	cardbox, _ := model.GetCardboxByID(boxID)
	pd.ActiveCardbox, _ = model.GetCardbox(cardbox[0])
	if len(pd.ActiveCardbox.Cards) == 0 {
		pd.ActiveCards = make([]model.Card, 0)

		pd.ActiveCard = model.Card{}

	} else {
		pd.ActiveCards = pd.ActiveCardbox.Cards
		pd.ActiveCard = pd.ActiveCards[0]

		userMap, _ := model.GetUser(pd.Userinfo.Username)
		user := model.GetUserDoc(userMap)
		userCardboxes := user.Cardboxes
		for _, i := range userCardboxes {
			if i.BoxID == pd.ActiveCardbox.Id {
				pd.ActiveCardbox.Progress = i.Progress

				lvls := make([]model.LevelMeter, 5)
				lvls[0].Level = "0"
				lvls[1].Level = "1"
				lvls[2].Level = "2"
				lvls[3].Level = "3"
				lvls[4].Level = "4"
				for i := 0; i < 5; i++ {
					lvls[i].LevelCode = "none"
				}

				if contains(i.Lvl0, pd.ActiveCard.CardID) {
					lvls[0].LevelCode = "currentLevel"
				}
				if contains(i.Lvl1, pd.ActiveCard.CardID) {
					lvls[1].LevelCode = "currentLevel"
				}
				if contains(i.Lvl2, pd.ActiveCard.CardID) {
					lvls[2].LevelCode = "currentLevel"
				}
				if contains(i.Lvl3, pd.ActiveCard.CardID) {
					lvls[3].LevelCode = "currentLevel"
				}
				if contains(i.Lvl4, pd.ActiveCard.CardID) {
					lvls[4].LevelCode = "currentLevel"
				}

				pd.ActiveCard.Levelmeter = lvls
			}
		}
	}
	session.Values["pagedata"] = pd
	session.Save(req, w)
	http.Redirect(w, req, "/learn-box", http.StatusSeeOther)
}

func RightAnswer(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "session")
	pd := session.Values["pagedata"].(pageData)
	user, _ := model.GetUser(pd.Userinfo.Username)
	userboxes := model.ToMapArray(user["cardboxes"])
	var progressBox map[string]interface{}
	for i := 0; i < len(userboxes); i++ {
		if userboxes[i]["boxID"] == pd.ActiveCardbox.Id {
			progressBox = userboxes[i]
		}
	}
	var newLvl string

	switch pd.Levels.ActiveCardLevel {
	case "lvl0":
		newLvl = "lvl1"
	case "lvl1":
		newLvl = "lvl2"
	case "lvl2":
		newLvl = "lvl3"
	case "lvl3":
		newLvl = "lvl4"
	case "lvl4":
		newLvl = "lvl4"
	}

	progressBox[newLvl] = append(progressBox[newLvl].([]interface{}), pd.ActiveCard.CardID)
	//var marker int
	for i := 0; i < len(progressBox[pd.Levels.ActiveCardLevel].([]interface{})); i++ {
		if progressBox[pd.Levels.ActiveCardLevel].([]interface{})[i] == pd.ActiveCard.CardID {
			progressBox[pd.Levels.ActiveCardLevel].([]interface{})[i] = progressBox[pd.Levels.ActiveCardLevel].([]interface{})[len(progressBox[pd.Levels.ActiveCardLevel].([]interface{}))-1]
			progressBox[pd.Levels.ActiveCardLevel] = progressBox[pd.Levels.ActiveCardLevel].([]interface{})[:len(progressBox[pd.Levels.ActiveCardLevel].([]interface{}))-1]
			break
		}
	}

	//	progressBox[pd.Levels.ActiveCardLevel] = append(progressBox[pd.Levels.ActiveCardLevel].([]interface{})[:marker], progressBox[pd.Levels.ActiveCardLevel].([]interface{})[marker+1:])

	for i := 0; i < len(userboxes); i++ {
		if userboxes[i]["boxID"] == pd.ActiveCardbox.Id {
			userboxes[i] = progressBox
		}
	}

	user["cardboxes"] = userboxes
	model.UpdateFile(user, user["_id"].(string))

	session.Values["pagedata"] = pd
	http.Redirect(w, req, "/learn/box/"+pd.ActiveCardbox.Id, http.StatusSeeOther)
}

func WrongAnswer(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "session")
	pd := session.Values["pagedata"].(pageData)
	user, _ := model.GetUser(pd.Userinfo.Username)
	userboxes := model.ToMapArray(user["cardboxes"])
	var progressBox map[string]interface{}
	for i := 0; i < len(userboxes); i++ {
		if userboxes[i]["boxID"] == pd.ActiveCardbox.Id {
			progressBox = userboxes[i]
		}
	}
	newLvl := "lvl0"

	progressBox[newLvl] = append(progressBox[newLvl].([]interface{}), pd.ActiveCard.CardID)
	//var marker int
	for i := 0; i < len(progressBox[pd.Levels.ActiveCardLevel].([]interface{})); i++ {
		if progressBox[pd.Levels.ActiveCardLevel].([]interface{})[i] == pd.ActiveCard.CardID {
			progressBox[pd.Levels.ActiveCardLevel].([]interface{})[i] = progressBox[pd.Levels.ActiveCardLevel].([]interface{})[len(progressBox[pd.Levels.ActiveCardLevel].([]interface{}))-1]
			progressBox[pd.Levels.ActiveCardLevel] = progressBox[pd.Levels.ActiveCardLevel].([]interface{})[:len(progressBox[pd.Levels.ActiveCardLevel].([]interface{}))-1]
			break
		}
	}

	//	progressBox[pd.Levels.ActiveCardLevel] = append(progressBox[pd.Levels.ActiveCardLevel].([]interface{})[:marker], progressBox[pd.Levels.ActiveCardLevel].([]interface{})[marker+1:])

	for i := 0; i < len(userboxes); i++ {
		if userboxes[i]["boxID"] == pd.ActiveCardbox.Id {
			userboxes[i] = progressBox
		}
	}

	user["cardboxes"] = userboxes
	model.UpdateFile(user, user["_id"].(string))

	session.Values["pagedata"] = pd
	http.Redirect(w, req, "/learn/box/"+pd.ActiveCardbox.Id, http.StatusSeeOther)
}
