package main

import (
	"log"
	"net/http"

	"forum/database"
	"forum/server"
)

func main() {
	DB := database.IntDB()
	err := database.CreateTable(&DB)
	if err != nil {
		log.Fatal("Error creating tables: ", err)
		return
	}

	api := server.New(DB)

	http.HandleFunc("/singup", api.PageSingUp)
	http.HandleFunc("/singin", api.PageSingIn)

	http.HandleFunc("/user/singup", api.DB.SingUp)
	http.HandleFunc("/user/singin", api.DB.SingIn)

	http.HandleFunc("/js/", server.Server)
	http.HandleFunc("/style/", server.Server)

	err = http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("Error starting the server: ", err) // Handle error more effectively
	}
}
