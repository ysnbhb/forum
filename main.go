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
	http.HandleFunc("/logout", api.LogOut)
	http.HandleFunc("/user/signup", api.DB.SingUp)
	http.HandleFunc("/user/signin", api.DB.SingIn)
	http.HandleFunc("/user/exist", api.DB.Exist)
	http.HandleFunc("/post/lastId", api.DB.LastID)
	http.HandleFunc("/frontend/", server.Server)
	// api
	http.HandleFunc("/api/getCategorie", api.DB.GetCtg)
	http.HandleFunc("/api/posts", api.DB.GetPost)
	http.HandleFunc("/api/post", api.DB.OnePost)
	http.HandleFunc("/api/addPost", api.DB.AddPost)
	http.HandleFunc("/api/post/like", DB.LikePost)
	http.HandleFunc("/api/commant/like", DB.LikeCommat)
	http.HandleFunc("/api/commant", DB.GetCommat)

	// run server
	println("server is run in http://localhost:8081/")
	err = http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("Error starting the server: ", err)
	}
}
