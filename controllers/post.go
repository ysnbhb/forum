package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"forum/utils"
)

func (db *Date) GetPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	posts := []utils.Post{}
	query := `
    SELECT 
        user.user_name, 
        post.id, 
        post.title, 
        post.contant, 
        post.create_date
    FROM 
        user
    INNER JOIN 
        post 
    ON 
        post.user_id = user.id
    ORDER BY 
        post.id DESC
    LIMIT ? OFFSET ?
`

	row, err := db.DB.Query(query, 20, "0")
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer row.Close()

	for row.Next() {
		post := utils.Post{}
		if err := row.Scan(&post.UserName, &post.Id, &post.Title, &post.Contant, &post.Date); err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		// fmt.Println(post)
		posts = append(posts, post)
	}

	if err := row.Err(); err != nil {
		fmt.Println("Error during rows iteration:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(posts); err != nil {
		fmt.Println("Error encoding JSON:", err)
	}
}
