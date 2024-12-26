package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"forum/utils"
)

func (db *Date) FilterWithCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	category := r.URL.Query().Get("category")
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	limitInt := 20
	offsetInt := 0
	if l, err := strconv.Atoi(limit); err == nil {
		limitInt = l
	}
	if o, err := strconv.Atoi(offset); err == nil {
		offsetInt = o
	}

	if category == "" {
		http.Error(w, "category is required", http.StatusBadRequest)
		return
	}
	categoryID := db.GetIdCateg(category)
	query := `
        SELECT 
            user.user_name, 
            post.id, 
            post.title, 
            post.contant, 
            post.create_date, 
            post.img, 
            post.categories
        FROM 
            user 
        INNER JOIN 
            post 
        ON 
            post.user_id = user.id
        INNER JOIN 
            categories_post 
        ON 
            categories_post.post_id = post.id
        WHERE 
            categories_post.categorie_id = ? 
        ORDER BY 
            post.id DESC
        LIMIT ? OFFSET ?;`
	rows, err := db.DB.Query(query, categoryID, limitInt, offsetInt)
	if err != nil {
		fmt.Println("SQL Query Error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	posts := []utils.Post{}
	for rows.Next() {
		post := utils.Post{}
		categories := ""

		if err := rows.Scan(&post.UserName, &post.Id, &post.Title, &post.Contant, &post.Date, &post.ImgUrl, &categories); err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}

		post.Categories = strings.Split(categories, ",")
		post.Reaction = db.GetReaction(r, post.Id, "post_id")
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error during rows iteration:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if len(posts) == 0 {
		fmt.Println("No results found.")
		http.Error(w, "No results found", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(posts); err != nil {
		fmt.Println("Error encoding JSON:", err)
	}
}

func (db *Date) GetIdCateg(str string) int {
	categ := 0
	db.DB.QueryRow(`SELECT id WHERE name_categorie  = ?`, str).Scan(&categ)
	return categ
}
