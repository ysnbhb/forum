package main

import (
	"fmt"
	"net/http"

	"forum/database"
)

func main() {
	DB := database.Intalction()
	err := database.CreateTable(DB)
	if err != nil {
		fmt.Println(err)
		return
	}
	http.HandleFunc("/home", DB.LogUp)
	http.ListenAndServe(":8081", nil)
}
