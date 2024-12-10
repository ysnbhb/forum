package server

import (
	"html/template"
	"net/http"

	"forum/database"
	"forum/handul"
)

type Apiserve struct {
	DB   handul.Date
}

func (api *Apiserve) PageSingUp(w http.ResponseWriter, r *http.Request) {
	api.HandlePage(w, r, "singup.html")
}

func (api *Apiserve) PageSingIn(w http.ResponseWriter, r *http.Request) {
	api.HandlePage(w, r, "singin.html")
}

func (api *Apiserve) HandlePage(w http.ResponseWriter, r *http.Request, htmlfile string) {
	cookie, err := r.Cookie("token")
	if err == nil {
		// api.DB.CheckEXist(cookie.Value)
		db := database.IntDB()
		exist := db.CheckEXist(cookie.Value)
		if exist {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}
	tmp, err := template.ParseFiles("template/html/" + htmlfile)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	tmp.Execute(w, nil)
}

func New(DB handul.Date) *Apiserve {
	return &Apiserve{
		DB:   DB,
	}
}
