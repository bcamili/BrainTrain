package controller

import (
	"BrainTrain/app/model"
	"encoding/base64"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func Mycrds(w http.ResponseWriter, r *http.Request) {

	var pd pageData
	session, err := store.Get(r, "session")
	if session.Values["pagedata"] == nil {
		http.Redirect(w, r, "/index", http.StatusFound)
	} else {
		pd = session.Values["pagedata"].(pageData)

	}

	pd.IsLoggedIn = session.Values["authenticated"] == true

	var boxes []map[string]interface{}
	allUserBoxes, _ := model.GetAllCardboxesByUsername(pd.Userinfo.Username)
	one := map[string]interface{}{
		"categoryname": "Selbst erstellte Karteikarten",
		"Boxes":        allUserBoxes,
	}
	boxes = append(boxes, one)

	user, _ := model.GetUser(pd.Userinfo.Username)
	allSavedBoxes := model.ToMapArray(user["cardboxes"])
	//pd.MainInfo.UserBoxCount = len(allSavedBoxes)
	boxes2 := make([]map[string]interface{}, 0)
	for _, v := range allSavedBoxes {
		thisBoxID := v["boxID"].(string)

		thisBox, _ := model.GetCardboxByID(thisBoxID)
		if len(thisBox) == 0 {
			http.Redirect(w, r, "/delete/"+thisBoxID, http.StatusSeeOther)
			return
		} else {

			if thisBox[0]["author"] != pd.Userinfo.Username {
				boxes2 = append(boxes2, thisBox[0])
			}
		}
	}

	two := map[string]interface{}{
		"categoryname": "Gelernte Karteien anderer Nutzer",
		"Boxes":        boxes2,
	}
	boxes = append(boxes, two)

	boxesStruct := make([]Category, len(boxes))

	for i := 0; i < len(boxes); i++ {
		boxesStruct[i] = MapToCategory(boxes[i])
	}

	for i := 0; i < len(boxesStruct); i++ {
		structBoxes := boxesStruct[i].Boxes
		for j := 0; j < len(structBoxes); j++ {
			for k := 0; k < len(allSavedBoxes); k++ {
				if allSavedBoxes[k]["boxID"].(string) == structBoxes[j].Id {
					structBoxes[j].Progress = allSavedBoxes[k]["progress"].(float64)
				}
			}
		}

	}

	pd.Categories = boxesStruct
	pd.Title = "BrainTrain: Karteien"
	pd.ErrorMsg = Error{}
	pd.MainInfo = StatsU(pd.Userinfo.Username)
	session.Values["pagedata"] = pd
	err = session.Save(r, w)
	err = tpl.ExecuteTemplate(w, "mycards.gohtml", pd)

	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func DeleteBox(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "session")
	pd := session.Values["pagedata"].(pageData)
	boxID := req.URL.EscapedPath()[8:len(req.URL.EscapedPath())]
	box, _ := model.GetCardboxByID(boxID)
	if len(box) == 0 {
		user, _ := model.GetUser(pd.Userinfo.Username)
		cardboxes := model.ToMapArray(user["cardboxes"])
		for i := 0; i < len(cardboxes); i++ {
			if cardboxes[i]["boxID"] == boxID {
				cardboxes[i] = cardboxes[len(cardboxes)-1]
				cardboxes = cardboxes[:len(cardboxes)-1]
			}
		}
		user["cardboxes"] = cardboxes
		model.UpdateFile(user, user["_id"].(string))
		http.Redirect(w, req, "/mycards", http.StatusSeeOther)
		return
	} else {
		if box[0]["author"] == pd.Userinfo.Username {
			user, _ := model.GetUser(pd.Userinfo.Username)
			boxnum := user["boxnum"].(float64)
			boxnum -= float64(1)
			user["boxnum"] = boxnum
			pd.Userinfo.Boxnum = boxnum
			cardnum := user["cardnum"].(float64)
			numcards := box[0]["numcards"].(float64)
			cardnum -= numcards

			user["cardnum"] = cardnum

			pd.Userinfo.Cardnum = cardnum
			model.DeleteFile(box[0]["_id"].(string))
			model.UpdateFile(user, user["_id"].(string))
			http.Redirect(w, req, "/mycards", http.StatusSeeOther)
		} else {
			user, _ := model.GetUser(pd.Userinfo.Username)
			cardboxes := model.ToMapArray(user["cardboxes"])
			for i := 0; i < len(cardboxes); i++ {
				if cardboxes[i]["boxID"] == boxID {
					cardboxes[i] = cardboxes[len(cardboxes)-1]
					cardboxes = cardboxes[:len(cardboxes)-1]
				}
			}
			user["cardboxes"] = cardboxes
			model.UpdateFile(user, user["_id"].(string))
			http.Redirect(w, req, "/mycards", http.StatusSeeOther)
		}
	}
}

func EditCardbox(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	pd := session.Values["pagedata"].(pageData)
	pd.Title = "Karte barbeiten"
	pd.Edit = true
	for i := 0; i < len(pd.ActiveCards); i++ {
		pd.ActiveCards[i].Count = (i + 1)

	}
	session.Values["pagedata"] = pd
	session.Save(r, w)
	err := tpl.ExecuteTemplate(w, "new-card-2.gohtml", pd)

	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func EditBoxSelector(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	pd := session.Values["pagedata"].(pageData)
	boxID := r.URL.EscapedPath()[10:len(r.URL.EscapedPath())]
	cardbox, _ := model.GetCardboxByID(boxID)
	pd.ActiveCardbox, _ = model.GetCardbox(cardbox[0])
	pd.ActiveCards = pd.ActiveCardbox.Cards
	if len(pd.ActiveCards) != 0 {
		pd.ActiveCard = pd.ActiveCards[0]
		pd.ActiveCards[0].Select = "selectedCard"
	} else {
		pd.ActiveCards = make([]model.Card, 0)
		pd.ActiveCard = model.Card{}
	}
	session.Values["pagedata"] = pd
	session.Save(r, w)
	http.Redirect(w, r, "/edit", http.StatusSeeOther)
}

func EditCardSelector(w http.ResponseWriter, r *http.Request) {
	var pd pageData
	session, _ := store.Get(r, "session")
	if session.Values["pagedata"] == nil {
		pd = pageData{}
	} else {
		pd = session.Values["pagedata"].(pageData)

	}
	cardID := r.URL.EscapedPath()[15:len(r.URL.EscapedPath())]
	for i := 0; i < len(pd.ActiveCards); i++ {
		pd.ActiveCards[i].Select = "none"
		if pd.ActiveCards[i].CardID == cardID {
			pd.ActiveCard = pd.ActiveCards[i]
			pd.ActiveCards[i].Select = "selectedCard"

			session, _ := store.Get(r, "session")

			if session.Values["authenticated"] != nil && session.Values["authenticated"] != false {
				pd.IsLoggedIn = session.Values["authenticated"].(bool)
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
		}

	}
	pd.Edit = true
	session.Values["pagedata"] = pd
	session.Save(r, w)
	http.Redirect(w, r, "/edit", http.StatusSeeOther)
}
func EditBoxInfo(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	pd := session.Values["pagedata"].(pageData)
	pd.Title = "Kartei bearbeiten"
	pd.Edit = true
	session.Values["pagedata"] = pd
	session.Save(r, w)
	err := tpl.ExecuteTemplate(w, "new-card.gohtml", pd)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func SaveCardbox(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	pd := session.Values["pagedata"].(pageData)
	if r.FormValue("category") == "empty" {
		pd.ErrorMsg.Error01 = true
		session.Values["pagedata"] = pd
		session.Save(r, w)
		http.Redirect(w, r, "/edit-info", http.StatusFound)

	} else {

		newTitle := r.FormValue("cardTitle")
		newDescription := r.FormValue("description")
		newVisibility := r.FormValue("visibility")
		newCategory, newSubcategory := ParseCategory(r.FormValue("category"))
		box, _ := model.GetCardboxByID(pd.ActiveCardbox.Id)
		box[0]["boxname"] = newTitle
		box[0]["description"] = newDescription
		box[0]["visibility"] = newVisibility
		box[0]["category"] = newCategory
		box[0]["subcategory"] = newSubcategory
		pd.ActiveCardbox, _ = model.GetCardbox(box[0])
		pd.ActiveCardbox.Id, _ = model.UpdateFile(box[0], box[0]["_id"].(string))
		session.Values["pagedata"] = pd
		session.Save(r, w)
		http.Redirect(w, r, "/edit", http.StatusSeeOther)
	}
}

func SaveCard(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	pd := session.Values["pagedata"].(pageData)

	newTitle := r.FormValue("cardTitle")
	newQuestion := r.FormValue("question")
	newAnswer := r.FormValue("answer")

	box, _ := model.GetCardboxByID(pd.ActiveCardbox.Id)
	var newCards []map[string]interface{}
	if len(pd.ActiveCards) == 0 {
		hashedID, _ := bcrypt.GenerateFromPassword([]byte(strconv.FormatInt(time.Now().UnixNano(), 5)), 14)
		b64HashedID := base64.StdEncoding.EncodeToString(hashedID)
		newCard := model.Card{
			CardID:   b64HashedID,
			Title:    newTitle,
			Question: newQuestion,
			Answer:   newAnswer,
		}
		pd.ActiveCard = newCard
		mapCard := model.ToMap(pd.ActiveCard)
		delete(mapCard, "Select")
		delete(mapCard, "Levelmeter")
		delete(mapCard, "Count")
		pd.ActiveCards = append(pd.ActiveCards, pd.ActiveCard)
		newCards = make([]map[string]interface{}, 1)
		newCards[0] = mapCard
		box[0]["cards"] = newCards
		box[0]["numcards"] = 1
		pd.ActiveCards[0].Select = "selected"
		pd.ActiveCardbox.NumCards += 1
		user, _ := model.GetUser(pd.Userinfo.Username)
		ucardnum := user["cardnum"].(float64)
		ucardnum += 1
		user["cardnum"] = ucardnum
		model.UpdateFile(user, user["_id"].(string))
		pd.Userinfo, _ = model.GetUserByUsername(pd.Userinfo.Username)
	} else {
		for i := 0; i < len(pd.ActiveCards); i++ {
			if pd.ActiveCards[i].CardID == pd.ActiveCard.CardID {
				pd.ActiveCards[i].Title = newTitle
				pd.ActiveCards[i].Question = newQuestion
				pd.ActiveCards[i].Answer = newAnswer
				pd.ActiveCard = pd.ActiveCards[i]
				newCards = model.ToMapArray(pd.ActiveCards)
				mapCard := model.ToMap(pd.ActiveCard)
				delete(mapCard, "Select")
				delete(mapCard, "Levelmeter")
				delete(mapCard, "Count")
				newCards[i] = mapCard
			}
		}
		box[0]["cards"] = newCards
	}
	pd.ActiveCardbox.Id, _ = model.UpdateFile(box[0], box[0]["_id"].(string))
	session.Values["pagedata"] = pd
	session.Save(r, w)
	http.Redirect(w, r, "/edit", http.StatusSeeOther)
}

func DeleteCard(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	pd := session.Values["pagedata"].(pageData)

	box, _ := model.GetCardboxByID(pd.ActiveCardbox.Id)
	cards := model.ToMapArray(box[0]["cards"])

	cardID := r.URL.EscapedPath()[13:len(r.URL.EscapedPath())]
	for i, v := range cards {
		if v["cardID"].(string) == cardID {
			v = cards[len(cards)-1]
			v["cardID"] = cardID
			pd.ActiveCards[i] = pd.ActiveCards[len(pd.ActiveCards)-1]
			pd.ActiveCards[i].CardID = cardID
			pd.ActiveCards[i].Select = pd.ActiveCard.Select
			pd.ActiveCard = pd.ActiveCards[i]
		}
	}
	newCards := cards[:len(cards)-1]
	pd.ActiveCards = pd.ActiveCards[:len(pd.ActiveCards)-1]
	box[0]["cards"] = newCards
	box[0]["numcards"] = len(newCards)
	id, _ := model.UpdateFile(box[0], box[0]["_id"].(string))
	carbox, _ := model.GetCardboxByID(id)
	pd.ActiveCardbox, _ = model.GetCardbox(carbox[0])
	user, _ := model.GetUser(pd.Userinfo.Username)
	ucardnum := user["cardnum"].(float64)
	ucardnum -= 1
	user["cardnum"] = ucardnum
	model.UpdateFile(user, user["_id"].(string))
	pd.Userinfo, _ = model.GetUserByUsername(pd.Userinfo.Username)
	session.Values["pagedata"] = pd
	session.Save(r, w)
	http.Redirect(w, r, "/edit", http.StatusSeeOther)
}

func CreateCard(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	pd := session.Values["pagedata"].(pageData)
	box, _ := model.GetCardboxByID(pd.ActiveCardbox.Id)
	hashedID, _ := bcrypt.GenerateFromPassword([]byte(strconv.FormatInt(time.Now().UnixNano(), 5)), 14)
	b64HashedID := base64.StdEncoding.EncodeToString(hashedID)
	pd.ActiveCard = model.Card{
		CardID: b64HashedID,
		Select: "selected",
	}

	//	hier ist index out of range
	pd.ActiveCards = append(pd.ActiveCards, pd.ActiveCard)
	box[0]["cards"] = pd.ActiveCards
	box[0]["numcards"] = len(pd.ActiveCards)
	user, _ := model.GetUser(pd.Userinfo.Username)
	ucardnum := user["cardnum"].(float64)
	ucardnum += 1
	user["cardnum"] = ucardnum
	model.UpdateFile(user, user["_id"].(string))
	pd.Userinfo, _ = model.GetUserByUsername(pd.Userinfo.Username)
	id, _ := model.UpdateFile(box[0], box[0]["_id"].(string))
	carbox, _ := model.GetCardboxByID(id)
	pd.ActiveCardbox, _ = model.GetCardbox(carbox[0])
	session.Values["pagedata"] = pd
	session.Save(r, w)
	http.Redirect(w, r, "/edit", http.StatusSeeOther)
}
