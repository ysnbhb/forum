package server

import (
	"net/http"
	"os"
	"strings"
	"time"

	"forum/utils"
)

func Server(w http.ResponseWriter, r *http.Request) {
	filename := "." + r.URL.Path
	file, err := os.ReadFile(filename)
	if err != nil {
		utils.ErrorHandler(w, http.StatusNotFound, "Page Not Found", "The page you are looking for is not available!", nil)
		return
	}
	http.ServeContent(w, r, filename, time.Now(), strings.NewReader(string(file)))
}
