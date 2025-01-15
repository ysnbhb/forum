package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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
    ORDER BY 
        post.id DESC
    LIMIT ? OFFSET ?
`
	row, err := db.DB.Query(query, limitInt, offsetInt)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	userId := db.GetIdFromReq(r)
	defer row.Close()

	for row.Next() {
		Categories := ""
		post := utils.Post{}
		if err := row.Scan(&post.UserName, &post.Id, &post.Title, &post.Contant, &post.Date, &post.ImgUrl, &Categories); err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		post.Categories = strings.Split(Categories, " ,")
		post.Reaction = db.GetReaction(userId, post.Id, "post_id")
		post.Commant = db.GetCommatBYPost(post.Id, userId)
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
        post.create_date, 
		post.categories ,
		post.img
    FROM 
        user
    INNER JOIN 
        post 
    ON 
        post.user_id = user.id
    WHERE post.id = ?
		`
	categories := ""
	usrId := db.GetIdFromReq(r)
	err = db.DB.QueryRow(query, postId).Scan(&post.UserName, &post.Id, &post.Title, &post.Contant, &post.Date, &categories, &post.ImgUrl)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": http.StatusText(http.StatusNotFound)})
		return
	}
	post.Categories = strings.Split(categories, " ,")
	post.Reaction = db.GetReaction(usrId, post.Id, "post_id")
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
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
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
	 	INSERT INTO post(user_id , title , contant , img , categories)
		VALUES (? ,? ,? , ? , ?)
	`
	post := utils.Post{}
	file, mut, err := r.FormFile("img")
	if err == nil {
		isimg := utils.IsImage(mut)
		if !isimg {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "file is not image"})
			return
		}
		if mut.Size > 20<<20 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "size of image is more than expect"})
			return
		}
		file, err := SaveImg(file, mut)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "can't save file"})
			return
		}
		post.ImgUrl = file
	}
	if db.HandulPost(r, &post) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "your input is not correct"})
		return
	}
	categories := strings.Join(post.Categories, " ,")
	res, err := db.DB.Exec(query, id, post.Title, post.Contant, post.ImgUrl, categories)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "post not save correct try next time"})
		return
	}
	Postid, _ := res.LastInsertId()
	post.UserName, _ = db.GetName(id)
	post.Date = time.Now().Format("2006-01-02 15:04:05")
	post.Id = int(Postid)
	db.SaveCategories(post.Id, post.Idscategories)
	json.NewEncoder(w).Encode(post)
}

func SaveImg(file multipart.File, filehedre *multipart.FileHeader) (string, error) {
	uidImg, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	nameFiel := filepath.Join("userImg", uidImg.String()+filepath.Ext(filehedre.Filename))
	saveFile, err := os.Create(nameFiel)
	if err != nil {
		return "", err
	}
	defer saveFile.Close()
	_, err = io.Copy(saveFile, file)
	return nameFiel, err
}
