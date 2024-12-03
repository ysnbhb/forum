package main

import (
	"fmt"
	"net/http"

	"forum/database"
	"forum/server"
)

func main() {
	DB := database.Intalction()
	err := database.CreateTable(DB)
	if err != nil {
		fmt.Println(err)
		return
	}
	http.HandleFunc("/", server.PageSingUp)
	http.HandleFunc("/js/", server.Server)
	http.HandleFunc("/style/", server.Server)
	http.HandleFunc("/user/singup", DB.SingUp)
	http.HandleFunc("/user/check", DB.CheckEXist)
	http.ListenAndServe(":8080", nil)
}
