package main

import (
	"fmt"
	"net/http"

	"forum/database"
	"forum/server"
)

func main() {
	DB := database.IntDB()
	err := database.CreateTable(DB)
	if err != nil {
		fmt.Println(err)
		return
	}
	http.HandleFunc("/singup", server.PageSingUp)
	http.HandleFunc("/js/", server.Server)
	http.HandleFunc("/style/", server.Server)
	http.HandleFunc("/user/singup", DB.SingUp)
	http.ListenAndServe(":8080", nil)
}
