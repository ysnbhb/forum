package server

import (
	"net/http"

	"forum/utils"
)

func Page(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		utils.ErrorHandler(w, http.StatusNotFound, "Page not Fount", "The page you are looking for is not available!", nil)
		return
	}
}
