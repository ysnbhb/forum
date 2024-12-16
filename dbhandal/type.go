package dbhandal

import (
	"database/sql"
	"net/http"

	"forum/utils"
)

type BD interface {
	SingUp(http.ResponseWriter, *http.Request)
	SingIn(http.ResponseWriter, *http.Request)
	Insert(utils.User) (int, error)
	Select(string, string) (int, error)
	CraeteSession(int, string) error
	Exist(http.ResponseWriter, *http.Request)
	TakeName(http.ResponseWriter, string) bool
}

type Date struct {
	DB *sql.DB
}
