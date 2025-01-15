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
	category := r.URL.Query().Get("filterby")
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
	userId := db.GetIdFromReq(r)
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
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	posts := []utils.Post{}
	for rows.Next() {
		post := utils.Post{}
		categories := ""

		if err := rows.Scan(&post.UserName, &post.Id, &post.Title, &post.Contant, &post.Date, &post.ImgUrl, &categories); err != nil {
			fmt.Println(err)
			continue
		}

		post.Categories = strings.Split(categories, ",")
		post.Reaction = db.GetReaction(userId, post.Id, "post_id")
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(posts); err != nil {
		fmt.Println("Error encoding JSON:", err)
	}
}

func (db *Date) GetIdCateg(str string) int {
	categ := 0
	db.DB.QueryRow(`SELECT id FROM categories WHERE name_categorie  = ?`, str).Scan(&categ)
	return categ
}

func (db *Date) FilterMyPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": http.StatusText(http.StatusMethodNotAllowed)})
		return
	}

	userIs := db.GetIdFromReq(r)
	if userIs < 1 {
		w.WriteHeader(http.StatusNonAuthoritativeInfo)
		json.NewEncoder(w).Encode(map[string]string{"error": http.StatusText(http.StatusNonAuthoritativeInfo)})
		return
	}
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
	posts := []utils.Post{}
	query := `
    SELECT 
        user.user_name, 
        post.id, 
        post.title, 
        post.contant, 
        post.create_date , 
		post.img ,
		post.categories
    FROM 
        user
    INNER JOIN 
        post 
    ON 
        post.user_id = user.id
	WHERE
		post.user_id=  ?
    ORDER BY 
        post.id DESC
    LIMIT ? OFFSET ?
`
	row, err := db.DB.Query(query, userIs, limitInt, offsetInt)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer row.Close()

	for row.Next() {
		Categories := ""
		post := utils.Post{}
		if err := row.Scan(&post.UserName, &post.Id, &post.Title, &post.Contant, &post.Date, &post.ImgUrl, &Categories); err != nil {
			fmt.Println(err)
			continue
		}
		post.Categories = strings.Split(Categories, " ,")
		post.Reaction = db.GetReaction(userIs, post.Id, "post_id")
		posts = append(posts, post)
	}

	if err := row.Err(); err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(posts); err != nil {
		fmt.Println("Error encoding JSON:", err)
	}
}

func (db *Date) AllPost(w http.ResponseWriter, r *http.Request) {
	filterby := r.URL.Query().Get("filterby")
	if filterby == "all" || filterby == "" {
		db.GetPost(w, r)
	} else if filterby == "mypost" {
		db.FilterMyPost(w, r)
	} else if filterby == "likedpost" {
		db.FilterLikedPost(w, r)
	} else {
		db.FilterWithCategory(w, r)
	}
}

func (db *Date) FilterLikedPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

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

	// if !utils.AllowRequest(cookie.Value) {
	// 	utils.ErrorHandler(w, http.StatusTooManyRequests, http.StatusText(http.StatusTooManyRequests), "too many request you send", nil)
	// }
	userIs := db.GetIdFromReq(r)
	if userIs < 1 {
		w.WriteHeader(http.StatusNonAuthoritativeInfo)
		json.NewEncoder(w).Encode(map[string]string{"error": http.StatusText(http.StatusNonAuthoritativeInfo)})
		return
	}
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
			reaction 
		ON 
			reaction.post_id = post.id
		WHERE 
			reaction.user_id = ? AND reaction.type = ?
		ORDER BY 
			post.id DESC
		LIMIT ? OFFSET ?;
	`
	rows, err := db.DB.Query(query, userIs, "likes", limitInt, offsetInt)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	posts := []utils.Post{}
	for rows.Next() {
		post := utils.Post{}
		categories := ""

		if err := rows.Scan(&post.UserName, &post.Id, &post.Title, &post.Contant, &post.Date, &post.ImgUrl, &categories); err != nil {
			fmt.Println(err)
			continue
		}

		post.Categories = strings.Split(categories, ",")
		post.Reaction = db.GetReaction(userIs, post.Id, "post_id")
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error during rows iteration:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		fmt.Println("Error encoding JSON:", err)
	}
}

func (db *Date) FilterCommatedPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

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
	userIs := db.GetIdFromReq(r)
	if userIs < 1 {
		w.WriteHeader(http.StatusNonAuthoritativeInfo)
		json.NewEncoder(w).Encode(map[string]string{"error": http.StatusText(http.StatusNonAuthoritativeInfo)})
		return
	}
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
			comment 
		ON 
			comment.post_id = post.id
		WHERE 
			comment.user_id = ?
		GROUP BY 
			post.id
		ORDER BY 
			post.id DESC
		LIMIT ? OFFSET ?;
	`
	rows, err := db.DB.Query(query, userIs, limitInt, offsetInt)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	posts := []utils.Post{}
	for rows.Next() {
		post := utils.Post{}
		categories := ""

		if err := rows.Scan(&post.UserName, &post.Id, &post.Title, &post.Contant, &post.Date, &post.ImgUrl, &categories); err != nil {
			fmt.Println(err)
			continue
		}
		post.Categories = strings.Split(categories, " ,")
		post.Reaction = db.GetReaction(userIs, post.Id, "post_id")
		post.Commant = db.GetCommatBYPost(post.Id, userIs)
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error during rows iteration:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		fmt.Println("Error encoding JSON:", err)
	}
}

func (db *Date) GetCommatBYPost(post_id int, userId int) []utils.Commant {
	// commet := []utils.Commant{}
	comments := []utils.Commant{}
	query := `
	SELECT 
    	user.user_name, 
    	comment.id, 
    	comment.contant ,
		comment.create_date
	FROM 
    		user
	INNER JOIN 
    		comment 
	ON 
    		user.id = comment.user_id
	WHERE 
    		comment.post_id = ?
		ORDER BY 
    	comment.id DESC
	LIMIT ? OFFSET ?;
	`
	row, err := (db.DB.Query(query, post_id, 2, 0))
	if err != nil {
		return comments
	}
	for row.Next() {
		comment := utils.Commant{}
		row.Scan(&comment.UserName, &comment.Id, &comment.Contant, &comment.Date)
		comment.Reaction = db.GetReaction(userId, comment.Id, "comment_id")
		comments = append(comments, comment)
	}
	return comments
}

func (db *Date) FilterReactionedPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

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

	userIs := db.GetIdFromReq(r)
	if userIs < 1 {
		w.WriteHeader(http.StatusNonAuthoritativeInfo)
		json.NewEncoder(w).Encode(map[string]string{"error": http.StatusText(http.StatusNonAuthoritativeInfo)})
		return
	}
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
			reaction 
		ON 
			reaction.post_id = post.id
		WHERE 
			reaction.user_id = ?
		ORDER BY 
			post.id DESC
		LIMIT ? OFFSET ?;
	`
	rows, err := db.DB.Query(query, userIs, limitInt, offsetInt)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	posts := []utils.Post{}
	for rows.Next() {
		post := utils.Post{}
		categories := ""

		if err := rows.Scan(&post.UserName, &post.Id, &post.Title, &post.Contant, &post.Date, &post.ImgUrl, &categories); err != nil {
			fmt.Println(err)
			continue
		}

		post.Categories = strings.Split(categories, ",")
		post.Reaction = db.GetReaction(userIs, post.Id, "post_id")
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error during rows iteration:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		fmt.Println("Error encoding JSON:", err)
	}
}
