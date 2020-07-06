package controller

import (
	"BrainTrain/app/model"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func Crdbx(w http.ResponseWriter, req *http.Request) {
	pd := pageData{}
	session, _ := store.Get(req, "session")
	if session.Values["pagedata"] == nil {
		pd = pageData{}
	} else {
		pd = session.Values["pagedata"].(pageData)
	}
	if session.Values["authenticated"] == true {
		pd.IsLoggedIn = true
		pd.MainInfo = StatsU(pd.Userinfo.Username)

	} else {
		pd.IsLoggedIn = false
		pd.MainInfo = Stats()

	}

	if !(req.FormValue("sortBy") == "" || req.FormValue("sortBy") == "empty" || req.FormValue("sortBy") == "all") {
		fmt.Println("FormValue: " + req.FormValue("sortBy"))
		cat, subcat := ParseCategory(req.FormValue("sortBy"))
		boxes, _ := model.GetSubcatCardboxes(cat, subcat)
		pd.Categories = make([]Category, 1)
		pd.Categories[0].CategoryName = subcat
		pd.Categories[0].Boxes = boxes
		pd.Title = "BrainTrain: Karteikasten"
		pd.ErrorMsg = Error{}
		session.Values["pagedata"] = pd
		session.Save(req, w)
		err := tpl.ExecuteTemplate(w, "cardbox.gohtml", pd)

		if err != nil {
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		return
	} else {

		boxes, err := model.GetPublicCardboxes()
		if err != nil {
			log.Println(err)
			http.Error(w, "No boxes to display", http.StatusInternalServerError)
			return
		}

		var Cats map[string]bool
		Cats = make(map[string]bool)

		for i := 0; i < len(boxes); i++ {
			if !Cats[boxes[i]["category"].(string)] {
				Cats[boxes[i]["category"].(string)] = true

			}
		}
		var data []map[string]interface{}

		for k := range Cats {
			var boxes2 []map[string]interface{}
			for i := 0; i < len(boxes); i++ {
				if boxes[i]["category"] == k {
					boxes2 = append(boxes2, boxes[i])
				}
			}
			var toAppend map[string]interface{}
			toAppend = make(map[string]interface{})

			toAppend["categoryname"] = k
			toAppend["Boxes"] = boxes2

			data = append(data, toAppend)
		}

		dataStruct := make([]Category, len(data))

		for i := 0; i < len(data); i++ {
			dataStruct[i] = MapToCategory(data[i])
		}
		pd.Categories = dataStruct
		pd.Title = "BrainTrain: Karteikasten"
		pd.ErrorMsg = Error{}
		session.Values["pagedata"] = pd
		session.Save(req, w)
		err = tpl.ExecuteTemplate(w, "cardbox.gohtml", pd)

		if err != nil {
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
