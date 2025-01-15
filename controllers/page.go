package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"forum/utils"

	"github.com/gofrs/uuid/v5"
)

func (db *Date) SingUp(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(r.Referer())
	if r.Method == http.MethodGet {
		utils.ErrorHandler(w, http.StatusNotFound, "Page not Fount", "The page you are looking for is not available!", nil)
		return
	} else if r.Method != http.MethodPost {
		w.Header().Set("Content-type", "application/json")
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
	statuscode, id, err := db.Insert(user, "sing")
	if err != nil {
		w.WriteHeader(statuscode)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	err = db.SetCookie(w, int(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "there are error in server try later please")
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "user saved."})
	delete(utils.RateLimitData, utils.Ip(r.RemoteAddr))
}

func (db *Date) SingIn(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		utils.ErrorHandler(w, http.StatusNotFound, "Page not Fount", "The page you are looking for is not available!", nil)
		return
	} else if r.Method != http.MethodPost {
		utils.ErrorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed), "this Page doesn't support you method", nil)
		return

	}
	userInf := r.FormValue("userInf")
	userInf = strings.TrimLeft(userInf, " ")
	passwd := r.FormValue("passwd")
	if !utils.IsValidEmail(userInf) && !utils.IsValidUsername(userInf) || passwd == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "check you input")
		return
	}
	id, err := db.Select(userInf, passwd)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%v", err)
		return
	}
	err = db.SetCookie(w, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "there are error in server try later please")
		return
	}
	delete(utils.RateLimitData, utils.Ip(r.RemoteAddr))
}

func (db *Date) Exist(w http.ResponseWriter, r *http.Request) {
	db.TakeName(w, r)
}

func (db *Date) GetCtg(w http.ResponseWriter, r *http.Request) {
	categories := []string{}
	query := `
		SELECT (name_categorie) FROM categories
	`
	row, err := db.DB.Query(query)
	if err != nil {
		fmt.Println(err)
		return
	}
	for row.Next() {
		var categorie string
		row.Scan(&categorie)
		categories = append(categories, categorie)
	}
	json.NewEncoder(w).Encode(categories)
}

func (db *Date) AuthName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "method not allown"})
		return
	}
	cookie, err := r.Cookie("auth")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": http.StatusText(http.StatusUnauthorized)})
		return
	}
	user := utils.User{}
	user.Email, err = db.TakeEmail(cookie.Value)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": http.StatusText(http.StatusUnauthorized)})
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "there are error in server try other time"})
		return
	}
	user.User_name = r.FormValue("name")
	if !utils.IsValidUsername(user.User_name) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "check Your input"})
		return
	}
	user.Passwd, err = utils.GenratePass()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "there are error in server try other time"})
		return
	}
	staus, id, err := db.Insert(user, "auth")
	if err != nil {
		w.WriteHeader(staus)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	db.DeleteNoLog(cookie.Value)
	err = db.SetCookie(w, int(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "there are error in server try later please")
		return
	}
}

func (db *Date) DeleteNoLog(token string) {
	query := `DELETE FROM noLog WHERE token = ?`
	db.DB.Exec(query, token)
}

func (db *Date) SetCookie(w http.ResponseWriter, id int) error {
	uid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	err = db.CraeteSession(int(id), uid.String())
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    uid.String(),
		MaxAge:   int(time.Hour),
		HttpOnly: true,
		Path:     "/",
	})
	return nil
}

func (db *Date) TakeEmail(str string) (string, error) {
	email := ""
	query := `SELECT email FROM noLog WHERE token = ?`
	err := db.DB.QueryRow(query, str).Scan(&email)
	return email, err
}
