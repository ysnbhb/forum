package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"forum/utils"
)

func (db *Date) HandulPost(r *http.Request, post *utils.Post) bool {
	post.Title = strings.TrimSpace(r.FormValue("title"))
	post.Contant = strings.TrimSpace(r.FormValue("content"))
	json.NewDecoder(strings.NewReader(r.FormValue("categories"))).Decode(&post.Categories)
	post.Idscategories = db.ValidCateg(post.Categories)
	return post.Title == "" || post.Contant == "" || post.Idscategories == nil
}

func (db *Date) ValidCateg(categories []string) []int {
	query := `
		SELECT id FROM categories WHERE name_categorie = ?
	`
	idscategories := []int{}
	if len(categories) == 0 {
		return nil
	}
	for _, categorie := range categories {
		id := 0
		err := db.DB.QueryRow(query, categorie).Scan(&id)
		if err != nil || id == 0 {
			return nil
		}
		idscategories = append(idscategories, id)
	}
	return idscategories
}
