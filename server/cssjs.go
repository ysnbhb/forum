package server

import (
	"net/http"
	"os"
	"strings"
	"time"
)

func Server(w http.ResponseWriter, r *http.Request) {
	filename := "." + r.URL.Path
	file, err := os.ReadFile(filename)
	if err != nil {
		http.Error(w, "Page not Found", http.StatusNotFound)
		return
	}
	if strings.HasPrefix(r.URL.Path, "js") {
		w.Header().Set("Contant-Type", "text/javascript")
	} else {
		w.Header().Set("Contant-Type", "text/css")
	}
	http.ServeContent(w, r, filename, time.Now(), strings.NewReader(string(file)))
}
