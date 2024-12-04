package server

import (
	"forum/database"
	"html/template"
	"net/http"
)

func PageSingUp(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("template/html/index.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	cookie, err := r.Cookie("tookn")
	db := database.IntDB()
	if err == nil && cookie.HttpOnly {
		if db.CheckEXist(cookie.Value) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}
	tmp.Execute(w, nil)
}
