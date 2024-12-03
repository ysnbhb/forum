package server

import (
	"fmt"
	"html/template"
	"net/http"
)

func PageSingUp(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Page not Found", http.StatusNotFound)
		return
	}
	tmp, err := template.ParseFiles("template/html/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	tmp.Execute(w, nil)
}
