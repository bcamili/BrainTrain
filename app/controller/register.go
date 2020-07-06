package controller

import (
	"BrainTrain/app/model"
	"html/template"
	"log"
	"net/http"
)

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func Reg(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if session.Values["pagedata"] == nil {
		http.Redirect(w, r, "/index", http.StatusFound)
	} else {
		pd := session.Values["pagedata"].(pageData)
		pd.Title = "BrainTrain: Register"
		pd.IsLoggedIn = false
		pd.MainInfo = Stats()
		session.Values["pagedata"] = pd
		session.Save(r, w)
		err = tpl.ExecuteTemplate(w, "register.gohtml", pd)
		newError := Error{}
		pd.ErrorMsg = newError
		pd.MainInfo = Stats()
		session.Values["pagedata"] = pd
		session.Save(r, w)

		if err != nil {
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	checkPassword := r.FormValue("psw-repeat")
	tos := r.FormValue("tos")
	if password != checkPassword {
		session, err := store.Get(r, "session")
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		} else {
			pd := session.Values["pagedata"].(pageData)
			pd.ErrorMsg = Error{}
			pd.ErrorMsg.Error03 = true
			session.Values["pagedata"] = pd
			session.Save(r, w)
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}
	}

	allUsers, _, _ := model.GetAllUsers()
	for i := 0; i < len(allUsers); i++ {

		if allUsers[i]["email"].(string) == email {
			session, err := store.Get(r, "session")
			if err != nil {
				log.Println(err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			} else {
				pd := session.Values["pagedata"].(pageData)
				pd.ErrorMsg = Error{}
				pd.ErrorMsg.Error02 = true

				session.Values["pagedata"] = pd
				session.Save(r, w)
				http.Redirect(w, r, "/register", http.StatusSeeOther)
				return
			}
		}
	}

	if tos != "agreed" {
		session, err := store.Get(r, "session")
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		} else {
			pd := session.Values["pagedata"].(pageData)
			pd.ErrorMsg = Error{}
			pd.ErrorMsg.Error04 = true
			session.Values["pagedata"] = pd
			session.Save(r, w)
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}
	}

	user := model.User{}
	user.Username = username
	user.Password = password
	user.Email = email

	err := user.Add()
	user, err = model.GetUserByUsername(user.Username)

	if err != nil {

		session, err := store.Get(r, "session")
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		} else {
			pd := session.Values["pagedata"].(pageData)
			pd.ErrorMsg = Error{}
			pd.ErrorMsg.Error01 = true

			session.Values["pagedata"] = pd
			session.Save(r, w)
			http.Redirect(w, r, "/register", http.StatusSeeOther)
		}
	} else {
		session, _ := store.Get(r, "session")

		// Set user as authenticated
		session.Values["authenticated"] = true
		session.Values["username"] = username
		session.Options.MaxAge = 3600 * 3
		pd := session.Values["pagedata"].(pageData)

		allUserBoxes, _ := model.GetAllCardboxesByUsername(user.Username)

		cardnum := float64(0)
		for i := 0; i < len(allUserBoxes); i++ {
			cardnum += allUserBoxes[i]["numcards"].(float64)
		}

		pd.Userinfo = model.User{
			Username: user.Username,
			Email:    user.Email,
			Date:     user.Date,
			Boxnum:   user.Boxnum,
			Cardnum:  user.Cardnum,
			PicID:    user.PicID,
		}
		pd.MainInfo = Stats()
		allSavedBoxes := model.ToMapArray(user.Cardboxes)
		pd.MainInfo.UserBoxCount = len(allSavedBoxes)
		pd.IsLoggedIn = session.Values["authenticated"].(bool)
		session.Values["pagedata"] = pd
		err = session.Save(r, w)

		http.Redirect(w, r, "/mycards", http.StatusFound)
	}

}
