package server

import (
	"fmt"
	"html/template"
	"net/http"
)

func PageSingUp(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("template/html/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	tmp.Execute(w, nil)
}
