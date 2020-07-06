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

func Nwcrd(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "session")
	pd := session.Values["pagedata"].(pageData)

	pd.Title = "BrainTrain: Neue Kartei"

	session.Values["pagedata"] = pd
	session.Save(req, w)

	err = tpl.ExecuteTemplate(w, "new-card.gohtml", pd)
	pd.ErrorMsg = Error{}
	pd.MainInfo = StatsU(pd.Userinfo.Username)
	pd.Edit = false
	session.Values["pagedata"] = pd
	session.Save(req, w)

	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func Nwcrd2(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "session")
	pd := session.Values["pagedata"].(pageData)
	pd.Title = "BrainTrain: Neue Kartei"
	pd.MainInfo = StatsU(pd.Userinfo.Username)
	pd.Edit = false
	for i := 0; i < len(pd.ActiveCards); i++ {
		pd.ActiveCards[i].Count = (i + 1)

	}
	session.Values["pagedata"] = pd
	session.Save(req, w)
	err := tpl.ExecuteTemplate(w, "new-card-2.gohtml", pd)

	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

}
func Sbmtcrdbx(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	pd := session.Values["pagedata"].(pageData)
	pd.ErrorMsg = Error{}
	if r.FormValue("category") == "empty" {
		pd.ErrorMsg.Error01 = true
		session.Values["pagedata"] = pd
		session.Save(r, w)
		http.Redirect(w, r, "/new-card", http.StatusFound)

	} else {

		boxname := r.FormValue("cardTitle")
		category, subcategory := ParseCategory(r.FormValue("category"))
		description := r.FormValue("description")
		visibility := r.FormValue("visibility")

		newBox := model.Cardbox{
			Boxname:     boxname,
			Author:      pd.Userinfo.Username,
			Category:    category,
			Subcategory: subcategory,
			Description: description,
			Visibility:  visibility,
		}

		newID, err := newBox.Add()

		if err != nil {
			http.Redirect(w, r, "/new-card-2", http.StatusSeeOther)

		} else {
			user, _ := model.GetUser(pd.Userinfo.Username)
			boxnum := user["boxnum"].(float64)
			boxnum += float64(1)
			user["boxnum"] = boxnum
			pd.Userinfo.Boxnum = boxnum
			model.UpdateFile(user, user["_id"].(string))
			pd.ActiveCardbox = newBox
			pd.ActiveCardbox.Id = newID
			pd.ActiveCards = pd.ActiveCardbox.Cards
			session.Values["pagedata"] = pd
			session.Save(r, w)
			http.Redirect(w, r, "/edit-box/"+pd.ActiveCardbox.Id, http.StatusSeeOther)
		}
	}
}

func Sbmtcrd(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	pd := session.Values["pagedata"].(pageData)
	oldBox, _ := model.GetCardboxByID(pd.ActiveCardbox.Id)

	oldCards := model.ToMapArray(oldBox[0]["cards"])

	title := r.FormValue("cardTitle")
	question := r.FormValue("question")
	answer := r.FormValue("answer")

	hashedID, _ := bcrypt.GenerateFromPassword([]byte(strconv.FormatInt(time.Now().UnixNano(), 5)), 14)
	b64HashedID := base64.StdEncoding.EncodeToString(hashedID)

	newCard := model.Card{
		CardID:   b64HashedID,
		Title:    title,
		Question: question,
		Answer:   answer,
	}

	// pd.ActiveCardbox.Cards = append(pd.ActiveCardbox.Cards, newCard)
	// pd.ActiveCardbox.NumCards++
	// box := model.ToMap(pd.ActiveCardbox)

	card := model.ToMap(newCard)
	delete(card, "_id")
	delete(card, "_rev")
	delete(card, "Levelmeter")
	delete(card, "Select")

	oldCards = append(oldCards, card)
	oldBox[0]["cards"] = oldCards
	i := oldBox[0]["numcards"].(float64)
	oldBox[0]["numcards"] = i + float64(1)
	newID, _ := model.UpdateFile(oldBox[0], oldBox[0]["_id"].(string))
	user, _ := model.GetUser(pd.Userinfo.Username)
	cardnum := user["cardnum"].(float64)
	cardnum += float64(1)
	user["cardnum"] = cardnum
	pd.Userinfo.Cardnum = cardnum
	model.UpdateFile(user, user["_id"].(string))
	n, _ := model.GetCardboxByID(newID)
	pd.ActiveCardbox, _ = model.GetCardbox(n[0])
	pd.ActiveCards = pd.ActiveCardbox.Cards
	session.Values["pagedata"] = pd
	session.Save(r, w)
	http.Redirect(w, r, "/new-card-2", http.StatusSeeOther)

}
