package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"forum/utils"

	"github.com/gofrs/uuid/v5"
)

func (db *Date) GetPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": http.StatusText(http.StatusMethodNotAllowed)})
		return
	}
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

func (db *Date) OnePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": http.StatusText(http.StatusMethodNotAllowed)})
		return
	}
	postId, err := strconv.Atoi(r.URL.Query().Get("postid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": http.StatusText(http.StatusBadRequest)})
		return
	}
	post := utils.Post{}
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
    WHERE post.id = ?
		`
	err = db.DB.QueryRow(query, postId).Scan(&post.UserName, &post.Id, &post.Title, &post.Contant, &post.Date)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": http.StatusText(http.StatusNotFound)})
		return
	}
	if err := json.NewEncoder(w).Encode(post); err != nil {
		fmt.Println("Error encoding JSON:", err)
	}
}

func (db *Date) AddPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": http.StatusText(http.StatusMethodNotAllowed)})
		return
	}
	cookie, err := r.Cookie("token")
	if err != nil {
		w.WriteHeader(http.StatusNonAuthoritativeInfo)
		json.NewEncoder(w).Encode(map[string]string{"error": http.StatusText(http.StatusNonAuthoritativeInfo)})
		return
	}
	id := db.TakeId(cookie.Value)
	if id < 1 {
		w.WriteHeader(http.StatusNonAuthoritativeInfo)
		json.NewEncoder(w).Encode(map[string]string{"error": http.StatusText(http.StatusNonAuthoritativeInfo)})
		return
	}
	query := `
	 	INSERT INTO post(user_id , title , contant , img)
	 `
	post := utils.Post{}
	_ = post
	file, _, err := r.FormFile("img")
	if err == nil {
		isimg, err := utils.IsImage(file)
		if !isimg || err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "file is not image"})
			return
		}
		file, err := SaveImg(file)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "can't save file"})
			return
		}
		post.ImgUrl = file
	}
	if HandulPost(r, &post) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "your input is not correct"})
		return
	}
	_, err = db.DB.Exec(query, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "post not save correct try next time"})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"error": "post save good"})
}

func SaveImg(file io.Reader) (string, error) {
	uidImg, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	nameFiel := filepath.Join("userImg", uidImg.String())
	saveFiel, err := os.Create(nameFiel)
	if err != nil {
		return "", err
	}
	io.Copy(saveFiel, file)
	return nameFiel, nil
}

func HandulPost(r *http.Request, post *utils.Post) bool {
	post.Title = strings.TrimSpace(r.FormValue("title"))
	post.Contant = strings.TrimSpace(r.FormValue("contant"))
	json.NewDecoder(strings.NewReader(r.FormValue("categories"))).Decode(&post.Categories)
	return post.Title == "" || post.Contant == "" || len(post.Categories) == 0
}
