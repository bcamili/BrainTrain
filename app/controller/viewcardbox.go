package controller

import (
	"BrainTrain/app/model"
	"html/template"
	"log"
	"net/http"

	"github.com/russross/blackfriday"
)

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func Vwcrdbx(w http.ResponseWriter, req *http.Request) {
	var pd pageData
	session, err := store.Get(req, "session")
	if session.Values["pagedata"] == nil {
		pd = pageData{}
	} else {
		pd = session.Values["pagedata"].(pageData)

	}

	pd.Title = "BrainTrain: Kartei anschauen"
	if session.Values["authenticated"] == true {
		pd.IsLoggedIn = true
		pd.MainInfo = StatsU(pd.Userinfo.Username)

	} else {
		pd.IsLoggedIn = false
		pd.MainInfo = Stats()

	}
	for i := 0; i < len(pd.ActiveCards); i++ {
		pd.ActiveCards[i].Count = (i + 1)

	}
	pd.ActiveCard.QuestionHTML = template.HTML(blackfriday.MarkdownCommon([]byte(pd.ActiveCard.Question)))
	pd.ActiveCard.AnswerHTML = template.HTML(blackfriday.MarkdownCommon([]byte(pd.ActiveCard.Answer)))
	session.Values["pagedata"] = pd
	session.Save(req, w)
	err = tpl.ExecuteTemplate(w, "view-cardbox.gohtml", pd)

	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func Crdhndlr(w http.ResponseWriter, req *http.Request) {
	var pd pageData
	session, _ := store.Get(req, "session")
	if session.Values["pagedata"] == nil {
		pd = pageData{}
	} else {
		pd = session.Values["pagedata"].(pageData)

	}
	boxID := req.URL.EscapedPath()[6:len(req.URL.EscapedPath())]
	cardbox, _ := model.GetCardboxByID(boxID)
	pd.ActiveCardbox, _ = model.GetCardbox(cardbox[0])
	if len(pd.ActiveCardbox.Cards) == 0 {
		pd.ActiveCards = make([]model.Card, 0)

		pd.ActiveCard = model.Card{}

	} else {
		pd.ActiveCards = pd.ActiveCardbox.Cards
		pd.ActiveCard = pd.ActiveCards[0]
		pd.ActiveCards[0].Select = "selectedCard"

		if session.Values["authenticated"] == true {
			pd.IsLoggedIn = true
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

		} else {
			pd.IsLoggedIn = false
		}
	}
	session.Values["pagedata"] = pd
	session.Save(req, w)
	http.Redirect(w, req, "/view-cardbox", http.StatusSeeOther)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func Cardchooser(w http.ResponseWriter, req *http.Request) {
	var pd pageData
	session, _ := store.Get(req, "session")
	if session.Values["pagedata"] == nil {
		pd = pageData{}
	} else {
		pd = session.Values["pagedata"].(pageData)

	}
	cardID := req.URL.EscapedPath()[15:len(req.URL.EscapedPath())]
	for i := 0; i < len(pd.ActiveCards); i++ {
		pd.ActiveCards[i].Select = "none"
		if pd.ActiveCards[i].CardID == cardID {
			pd.ActiveCard = pd.ActiveCards[i]
			pd.ActiveCards[i].Select = "selectedCard"

			session, _ := store.Get(req, "session")

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

	pd.ActiveCard.QuestionHTML = template.HTML(blackfriday.MarkdownCommon([]byte(pd.ActiveCard.Question)))
	pd.ActiveCard.AnswerHTML = template.HTML(blackfriday.MarkdownCommon([]byte(pd.ActiveCard.Answer)))
	session.Values["pagedata"] = pd
	session.Save(req, w)
	http.Redirect(w, req, "/view-cardbox", http.StatusSeeOther)
}
