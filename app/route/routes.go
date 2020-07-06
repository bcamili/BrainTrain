package routes

import (
	"BrainTrain/app/controller"
	"net/http"

	"github.com/gorilla/mux"
)

func Routes() {
	r := mux.NewRouter()
	r.HandleFunc("/", controller.Idx).Methods("GET")
	r.HandleFunc("/index", controller.Idx).Methods("GET")
	r.HandleFunc("/login", controller.AuthenticateUser).Methods("POST")
	r.HandleFunc("/logout", controller.Logout).Methods("GET")
	r.HandleFunc("/register", controller.Reg).Methods("GET")
	r.HandleFunc("/signup", controller.AddUser).Methods("POST")
	r.HandleFunc("/cardbox", controller.Crdbx).Methods("GET")
	r.HandleFunc("/mycards", controller.Auth(controller.Mycrds)).Methods("GET")
	r.HandleFunc("/delete/{{.Id}}", controller.Auth(controller.DeleteBox)).Methods("GET")
	r.HandleFunc("/delete-card/{{.CardID}}", controller.Auth(controller.DeleteCard)).Methods("GET")
	r.HandleFunc("/edit", controller.Auth(controller.EditCardbox)).Methods("GET")
	r.HandleFunc("/edit-info", controller.Auth(controller.EditBoxInfo)).Methods("GET")
	r.HandleFunc("/edit-box/{{.Id}}", controller.Auth(controller.EditBoxSelector)).Methods("GET")
	r.HandleFunc("/edit-selected/{{.CardID}}", controller.EditCardSelector).Methods("GET")
	r.HandleFunc("/create-card", controller.Auth(controller.CreateCard)).Methods("GET")
	r.HandleFunc("/save-card", controller.Auth(controller.SaveCard)).Methods("POST")
	r.HandleFunc("/save-cardbox", controller.Auth(controller.SaveCardbox)).Methods("POST")
	r.HandleFunc("/my-profile", controller.Auth(controller.Myprfl)).Methods("GET")
	r.HandleFunc("/update-profile", controller.Auth(controller.UpdateProfile)).Methods("POST")
	r.HandleFunc("/delete-user", controller.Auth(controller.DeleteUser)).Methods("GET")
	r.HandleFunc("/new-card", controller.Auth(controller.Nwcrd)).Methods("GET")
	r.HandleFunc("/new-card-2", controller.Auth(controller.Nwcrd2)).Methods("GET")
	r.HandleFunc("/view-cardbox", controller.Vwcrdbx).Methods("GET")
	r.HandleFunc("/view/{{._id}}", controller.Crdhndlr).Methods("GET")
	r.HandleFunc("/view/selected/{{.CardID}}", controller.Cardchooser).Methods("GET")
	r.HandleFunc("/learn/box/{{._id}}", controller.Auth(controller.LearnBox)).Methods("GET")
	r.HandleFunc("/learn-box", controller.Auth(controller.Lrn)).Methods("GET")
	r.HandleFunc("/learn-2", controller.Auth(controller.Lrn2)).Methods("GET")
	r.HandleFunc("/correct", controller.Auth(controller.RightAnswer)).Methods("GET")
	r.HandleFunc("/false", controller.Auth(controller.WrongAnswer)).Methods("GET")
	r.HandleFunc("/skip", controller.Auth(controller.Lrn)).Methods("GET")
	r.HandleFunc("/submit-cardbox", controller.Auth(controller.Sbmtcrdbx)).Methods("POST")
	r.HandleFunc("/submit-card", controller.Auth(controller.Sbmtcrd)).Methods("POST")
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("static"))))
	server := http.Server{
		Addr:    ":80",
		Handler: r,
	}

	server.ListenAndServe()
}
