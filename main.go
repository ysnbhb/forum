package main

import (
	"log"
	"net/http"

	"forum/database"
	"forum/server"
)

func main() {
	DB := database.IntDB()
	err := database.CreateTable(DB)
	if err != nil {
		log.Fatal("Error creating tables: ", err)
		return
	}

	api := server.New(DB)

	http.HandleFunc("/signup", api.PageSingUp)
	http.HandleFunc("/", server.Page)
	http.HandleFunc("/signin", api.PageSingIn)
	http.HandleFunc("/logout", server.LogOut)
	http.HandleFunc("/user/signup", api.DB.SingUp)
	http.HandleFunc("/user/signin", api.DB.SingIn)
	http.HandleFunc("/user/exist", api.DB.Exist)
	http.HandleFunc("/forntend/", server.Server)
	println("server is run in http://localhost:8081/")
	err = http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("Error starting the server: ", err) // Handle error more effectively
	}
}
