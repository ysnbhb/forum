package dbhandal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"forum/utils"

	"github.com/gofrs/uuid/v5"
)

func (db *Date) SingUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		utils.ErrorHandler(w, http.StatusNotFound, "Page not Fount", "The page you are looking for is not available!", nil)
		return
	} else if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "method not allown"})
		return
	}
	user := utils.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	w.Header().Set("Content-type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid input for logup"})
		return
	} else if !utils.IsValidEmail(user.Email) || !utils.IsValidUsername(user.User_name) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid input for logup"})
		return
	}
	if !utils.IsValidPassword(user.Passwd) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "password is weak"})
		return
	}
	user.Passwd, err = utils.HasPassowd(user.Passwd)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "sorry but there are error in server try anther time"})
		return
	}
	statuscode, err := db.Insert(user)
	if err != nil {
		w.WriteHeader(statuscode)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "user saved."})
}

func (db *Date) SingIn(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		return
	} else if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return

	}
	userInf := r.FormValue("userInf")
	userInf = strings.TrimLeft(userInf, " ")
	passwd := r.FormValue("passwd")
	if !utils.IsValidEmail(userInf) && !utils.IsValidUsername(userInf) || passwd == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "check you input")
		fmt.Println("check you input")
		return
	}
	id, err := db.Select(userInf, passwd)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%v", err)
		return
	}
	uid, err := uuid.NewV4()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "there are error in server try later please")
		return
	}
	err = db.CraeteSession(id, uid.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "there are error in server try later please")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    uid.String(),
		MaxAge:   300,
		HttpOnly: true,
		Path:     "/",
	})
}

func (db *Date) Exist(w http.ResponseWriter, r *http.Request) {
	cookis, err := r.Cookie("token")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	db.TakeName(w, cookis.Value)
}
