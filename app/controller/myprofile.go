package controller

import (
	"BrainTrain/app/model"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func Myprfl(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	pd := session.Values["pagedata"].(pageData)
	pd.Title = "BrainTrain: Mein Profil"
	if session.Values["authenticated"] != nil {
		pd.IsLoggedIn = session.Values["authenticated"].(bool)
	} else {
		pd.IsLoggedIn = false
	}
	pd.MainInfo = StatsU(pd.Userinfo.Username)
	user, _ := model.GetUserByUsername(pd.Userinfo.Username)
	pd.Userinfo = user
	session.Values["pagedata"] = pd
	session.Save(r, w)
	err := tpl.ExecuteTemplate(w, "my-profile.gohtml", pd)
	newError := Error{}
	pd.ErrorMsg = newError
	session.Values["pagedata"] = pd
	session.Save(r, w)

	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	pd := session.Values["pagedata"].(pageData)
	user, err := model.GetUser(pd.Userinfo.Username)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	allUsers, _, err := model.GetAllUsers()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	newMail := r.FormValue("newMail")
	oldPW := r.FormValue("oldPW")
	newPW := r.FormValue("newPW")
	newPWcheck := r.FormValue("confirmPW")

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	passwordDB, _ := base64.StdEncoding.DecodeString(user["password"].(string))
	err = bcrypt.CompareHashAndPassword(passwordDB, []byte(oldPW))
	if err != nil {
		pd := session.Values["pagedata"].(pageData)
		pd.ErrorMsg = Error{}
		pd.ErrorMsg.Error02 = true
		session.Values["pagedata"] = pd
		session.Save(r, w)
		http.Redirect(w, r, "/my-profile", http.StatusSeeOther)
		return
	}

	if newMail != "" {
		for _, v := range allUsers {
			if v["email"].(string) == newMail {
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
					http.Redirect(w, r, "/my-profile", http.StatusSeeOther)
					return
				}
			}
		}
	}

	if newPW != "empty" {
		err = bcrypt.CompareHashAndPassword(passwordDB, []byte(newPW))
		if err == nil {
			pd := session.Values["pagedata"].(pageData)
			pd.ErrorMsg = Error{}
			pd.ErrorMsg.Error02 = true
			session.Values["pagedata"] = pd
			session.Save(r, w)
			http.Redirect(w, r, "/my-profile", http.StatusSeeOther)
			return
		}
	}

	if newPW != newPWcheck {
		pd := session.Values["pagedata"].(pageData)
		pd.ErrorMsg = Error{}
		pd.ErrorMsg.Error02 = true
		session.Values["pagedata"] = pd
		session.Save(r, w)
		http.Redirect(w, r, "/my-profile", http.StatusSeeOther)
		return
	}

	if newMail != "" {
		user["email"] = newMail
		pd.Userinfo.Email = newMail
	}

	if newPW != "" {
		hashedPwd, _ := bcrypt.GenerateFromPassword([]byte(newPW), 14)
		b64HashedPwd := base64.StdEncoding.EncodeToString(hashedPwd)
		user["password"] = b64HashedPwd
	}

	r.ParseMultipartForm(32 << 20)

	file, _, err := r.FormFile("profilePic")
	if err == nil {

		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		f, err := os.OpenFile("static/docs/"+pd.Userinfo.PicID, os.O_WRONLY|os.O_CREATE, 0666)

		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		_, err = io.Copy(f, file)
	}

	model.UpdateFile(user, user["_id"].(string))
	session.Values["pagedata"] = pd
	session.Save(r, w)
	http.Redirect(w, r, "/my-profile", http.StatusSeeOther)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	pd := session.Values["pagedata"].(pageData)
	cardboxes, _ := model.GetAllCardboxesByUsername(pd.Userinfo.Username)
	for _, box := range cardboxes {
		model.DeleteFile(box["_id"].(string))
	}
	model.DeleteFile(pd.Userinfo.ID)

	http.Redirect(w, r, "/logout", http.StatusSeeOther)
}
