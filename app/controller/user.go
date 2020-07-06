package controller

import (
	"BrainTrain/app/model"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"log"
	"math"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var store *sessions.FilesystemStore

func init() {

	key := make([]byte, 32)
	rand.Read(key)
	store = sessions.NewFilesystemStore(os.TempDir(), key)
	store.MaxLength(math.MaxInt64)
	gob.Register(pageData{})
}

func AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	var err error
	var user = model.User{}
	pd := pageData{}

	username := r.FormValue("username")
	password := r.FormValue("password")

	// Authentication
	user, err = model.GetUserByUsername(username)

	if err == nil {
		// decode base64 String to []byte
		passwordDB, _ := base64.StdEncoding.DecodeString(user.Password)
		err = bcrypt.CompareHashAndPassword(passwordDB, []byte(password))

		if err == nil {
			session, _ := store.Get(r, "session")

			// Set user as authenticated
			session.Values["authenticated"] = true
			session.Values["username"] = username
			session.Options.MaxAge = 3600 * 3

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
		} else {
			session, _ := store.Get(r, "session")
			pd.ErrorMsg.Error01 = true
			session.Values["pagedata"] = pd
			err = session.Save(r, w)
			http.Redirect(w, r, "/index", http.StatusFound)
		}
	} else {
		err = tpl.ExecuteTemplate(w, "register.gohtml", pd)

		if err != nil {
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

// Logout controller
func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	session.Values["authenticated"] = false
	session.Values["username"] = ""
	session.Values["pagedata"] = pageData{}
	session.Save(r, w)

	http.Redirect(w, r, "/index", http.StatusFound)
}

// Auth is an authentication handler
func Auth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		// Check if user is authenticated
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Redirect(w, r, "/index", http.StatusFound)
		} else {
			h(w, r)
		}
	}
}
