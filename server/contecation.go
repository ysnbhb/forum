package server

import (
	"net/http"

	"forum/controllers"
)

type BD interface {
	// server method of database
	SingUp(http.ResponseWriter, *http.Request)
	SingIn(http.ResponseWriter, *http.Request)
	Exist(http.ResponseWriter, *http.Request)
	LastID(http.ResponseWriter, *http.Request)
	GetCtg(http.ResponseWriter, *http.Request)
	GetPost(http.ResponseWriter, *http.Request)
	AddPost(http.ResponseWriter, *http.Request)
	OnePost(http.ResponseWriter, *http.Request)
	// some of method of insert or select
	DelectSeoin(string) error
	DeleteNoLog(string)
}

type Apiserve struct {
	DB BD
}

func New(DB *controllers.Date) *Apiserve {
	return &Apiserve{
		DB: DB,
	}
}
