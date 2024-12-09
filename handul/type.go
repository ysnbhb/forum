package handul

import (
	"database/sql"
	"net/http"
)

type BD interface {
	SingUp(http.ResponseWriter, *http.Request)
	SingIn(http.ResponseWriter, *http.Request)
	Insert(User) (int, error)
	Select(string, string) (int, error)
	CheckEXist(string) bool
	CraeteSession(int, string) error
}

type Date struct {
	DB *sql.DB
}

type User struct {
	User_name string `json:"user_name"`
	Email     string `json:"email"`
	Passwd    string `json:"passwd"`
}
