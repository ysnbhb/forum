package handul

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"unicode"

	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
)

func (db *Date) SingUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		return
	} else if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "method not allown"})
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
	if !IsValidUsername(user.User_name) || !IsValidEmail(user.Email) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "check you input , invalid input"})
		return
	}
	user.Passwd, err = HasPassowd(user.Passwd)
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
	json.NewEncoder(w).Encode(map[string]string{"status": "user insert into database"})
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
	if !IsValidEmail(userInf) && !IsValidUsername(userInf) || passwd == "" {
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

// This for valid username

func IsValidUsername(username string) bool {
	if username == "" {
		return false
	}
	fmt.Println(username)
	last := []rune(username)[0]
	for _, c := range username {
		if c == '_' {
			if last == '_' {
				return false
			}
			last = c
			continue
		}
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			return false
		}
		last = c
	}
	return true
}

// This for valid email

func IsValidEmail(email string) bool {
	re := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(re)
	return regex.MatchString(email)
}

func HasPassowd(password string) (string, error) {
	hashpassord, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashpassord), nil
}
