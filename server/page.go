package server

import (
	"html/template"
	"net/http"

	"forum/utils"
)

func Page(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		utils.ErrorHandler(w, http.StatusNotFound, "Page not Found", "The page you are looking for is not available!", nil)
		return
	}
	if r.Method != http.MethodGet {
		utils.ErrorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed), "this Page doesn't support your Method", nil)
		return
	}
	tmp, err := template.ParseFiles("view/home.html")
	if err != nil {
		utils.ErrorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "sorry but there are Error in server try next time", nil)
		return
	}
	tmp.Execute(w, nil)
}

func LogOut(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		HttpOnly: false,
		Value:    "",
		Name:     "token",
		MaxAge:   0,
		Path:     "/",
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
