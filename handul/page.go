package handul

import (
	"encoding/json"
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
	w.Header().Set("Content-type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid input for logup"})
		return
	} else if user.Email == "" || user.Passwd == "" || user.User_name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid input for logup"})
		return
	}
	err = db.Insert(user)
	if err != nil && (err.Error() == "email" || err.Error() == "user name") {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error() + " already used try anther " + err.Error()})
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "sorry but there are error in server try anther time"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "user insert into database"})

}

func (db *Date) CheckEXist(checker string) bool {
	exist := false
	err := db.DB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 
			FROM session
			WHERE uid = ?
		)`, checker, checker).Scan(&exist)
	return err == nil && exist
}
