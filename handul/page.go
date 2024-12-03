package handul

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (db *Date) SingUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		return
	} else if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "method not allown 1"})
		return
	}
	user := User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	fmt.Println(user)
	w.Header().Set("Content-type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid input for logup"})
		return
	} else if user.Email == "" || user.Passwd == "" || user.User_name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid input for logup"})
		return
	} else {
		err = db.Insert(user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "email or user name already exists"})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "user insert into database"})

	}
}

func (db *Date) CheckEXist(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	checker := r.FormValue("checker")
	exist := false
	err := db.DB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 
			FROM user
			WHERE user_name = ? OR email = ?
		)`, checker, checker).Scan(&exist)
	fmt.Println(exist, err)
	if err != nil || !exist {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(400)
	}
}
