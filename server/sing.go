package server

import (
	"html/template"
	"net/http"

	"forum/dbhandal"
)

type Apiserve struct {
	DB dbhandal.BD
}

func (api *Apiserve) PageSingUp(w http.ResponseWriter, r *http.Request) {
	api.HandlePage(w, r, "singup.html")
}

func (api *Apiserve) PageSingIn(w http.ResponseWriter, r *http.Request) {
	api.HandlePage(w, r, "singin.html")
}

func (api *Apiserve) HandlePage(w http.ResponseWriter, r *http.Request, htmlfile string) {
	tmp, err := template.ParseFiles("veiw/" + htmlfile)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	tmp.Execute(w, nil)
}

func New(DB *dbhandal.Date) *Apiserve {
	return &Apiserve{
		DB: DB,
	}
}
