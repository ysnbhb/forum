package server

import (
	"html/template"
	"net/http"

	"forum/controllers"
	"forum/utils"
)

type Apiserve struct {
	DB controllers.BD
}

func (api *Apiserve) PageSingUp(w http.ResponseWriter, r *http.Request) {
	api.HandlePage(w, r, "singup.html")
}

func (api *Apiserve) PageSingIn(w http.ResponseWriter, r *http.Request) {
	api.HandlePage(w, r, "singin.html")
}

func (api *Apiserve) HandlePage(w http.ResponseWriter, r *http.Request, htmlfile string) {
	if r.Method != http.MethodGet {
		utils.ErrorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed), "this Page doesn't support your Method", nil)
		return
	}
	tmp, err := template.ParseFiles("view/" + htmlfile)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	tmp.Execute(w, nil)
}

func New(DB *controllers.Date) *Apiserve {
	return &Apiserve{
		DB: DB,
	}
}
