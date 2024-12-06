package main

import (
	"fmt"
	"net/http"
	"time"

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
	http.HandleFunc("/singin", server.PageSingIn)
	http.HandleFunc("/js/", server.Server)
	http.HandleFunc("/style/", server.Server)
	http.HandleFunc("/user/singup", DB.SingUp)
	http.HandleFunc("/user/singin", DB.SingIn)
	tricker := time.NewTicker(time.Second * 2)
	go func() {
		for {
			DB.DeleteSession()
			<-tricker.C
		}
	}()
	http.ListenAndServe(":8080", nil)
}
