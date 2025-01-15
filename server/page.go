package server

import (
	"html/template"
	"net/http"

	"forum/utils"
)

func (api *Apiserve) Page(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		utils.ErrorHandler(w, http.StatusNotFound, "Page not Found", "The page you are looking for is not available!", nil)
		return
	}
	if r.Method != http.MethodGet {
		utils.ErrorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed), "this Page doesn't support your Method", nil)
		return
	}
	cookie, err := r.Cookie("auth")
	if err == nil {
		api.DB.DeleteNoLog(cookie.Value)
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			MaxAge: -1,
		})
	}
	tmp, err := template.ParseFiles("view/home.html")
	if err != nil {
		utils.ErrorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "sorry but there are Error in server try next time", nil)
		return
	}
	tmp.Execute(w, nil)
}

func AuthSing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ErrorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed), "this Page doesn't support your Method", nil)
		return
	}
	tmp, err := template.ParseFiles("view/auth.html")
	if err != nil {
		utils.ErrorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "sorry but there are Error in server try next time", nil)
		return
	}
	tmp.Execute(w, nil)
}

func (api *Apiserve) LogOut(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		utils.ErrorHandler(w, http.StatusNotFound, "Page not Found", "The page you are looking for is not available!", nil)
		return
	}
	err = api.DB.DelectSeoin(cookie.Value)
	if err != nil {
		utils.ErrorHandler(w, http.StatusNotFound, "Page not Found", "The page you are looking for is not available!", nil)
		return
	}
	delete(utils.RateLimitData, utils.Ip(r.RemoteAddr))
	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
