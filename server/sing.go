package server

import (
	"forum/database"
	"html/template"
	"net/http"
)

func PageSingUp(w http.ResponseWriter, r *http.Request) {
	HandlePage(w, r, "singup.html")
}

func PageSingIn(w http.ResponseWriter, r *http.Request) {
	HandlePage(w, r, "singin.html")
}

func HandlePage(w http.ResponseWriter, r *http.Request, htmlfile string) {
	cookie, err := r.Cookie("token")
	if err == nil {
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
