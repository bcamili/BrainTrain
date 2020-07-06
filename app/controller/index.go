package controller

import (
	"log"
	"net/http"
)

func Idx(w http.ResponseWriter, req *http.Request) {
	var pd pageData
	session, err := store.Get(req, "session")
	if session.Values["pagedata"] == nil {
		pd = pageData{MainInfo: stats}
	} else {
		pd = session.Values["pagedata"].(pageData)

	}

	if (err == nil) && (session.Values["authenticated"] == true) {
		pd.IsLoggedIn = true
		pd.MainInfo = StatsU(pd.Userinfo.Username)

	} else {
		pd.IsLoggedIn = false
		pd.MainInfo = Stats()

	}

	pd.Title = "BrainTrain: Start"
	if pd.ErrorMsg.Error01 != true {
		pd.ErrorMsg = Error{}
	}

	session.Values["pagedata"] = pd
	session.Save(req, w)

	err = tpl.ExecuteTemplate(w, "index.gohtml", pd)
	pd.ErrorMsg = Error{}
	session.Values["pagedata"] = pd
	session.Save(req, w)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
