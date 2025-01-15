package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"forum/utils"
)

func (db *Date) GetCommat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": http.StatusText(http.StatusMethodNotAllowed)})
		return
	}
	postid, err := strconv.Atoi(r.URL.Query().Get("postid"))
	if err != nil {
		http.Error(w, "post id is required", http.StatusBadRequest)
		return
	}
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	limitInt := 5
	offsetInt := 0
	if l, err := strconv.Atoi(limit); err == nil {
		limitInt = l
	}
	if o, err := strconv.Atoi(offset); err == nil {
		offsetInt = o
	}
	userId := db.GetIdFromReq(r)
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
	row, err := (db.DB.Query(query, postid, limitInt, offsetInt))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": http.StatusText(http.StatusBadRequest)})
		return
	}
	comments := []utils.Commant{}
	for row.Next() {
		comment := utils.Commant{}
		row.Scan(&comment.UserName, &comment.Id, &comment.Contant, &comment.Date)
		comment.Reaction = db.GetReaction(userId, comment.Id, "comment_id")
		comments = append(comments, comment)
	}
	json.NewEncoder(w).Encode(comments)
}

func (db *Date) AddCommat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": http.StatusText(http.StatusMethodNotAllowed)})
		return
	}
	cookie, err := r.Cookie("token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": http.StatusText(http.StatusUnauthorized)})
		return
	}
	id := db.TakeId(cookie.Value)
	if id < 1 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": http.StatusText(http.StatusUnauthorized)})
		return
	}
	postid, err := strconv.Atoi(r.FormValue("postid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid post id"})
		return
	}

	commant := utils.Commant{}
	commant.UserId = id
	exist := db.CheckCommat(postid)
	if exist != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "post not found"})
		return
	}
	commant.Post_id = postid
	commant.Contant = strings.TrimSpace(r.FormValue("contant"))
	fmt.Println(commant.Contant)
	if commant.Contant == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "no contant"})
		return
	}
	commant.UserName, _ = db.GetName(id)
	query := `
	 	INSERT INTO comment(user_id , post_id , contant)
		VALUES (? ,? ,?)
	`
	commant.Date = time.Now().Format("2006-01-02 15:04:05")
	res, err := db.DB.Exec(query, id, postid, commant.Contant)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}
	Postid, _ := res.LastInsertId()
	commant.Id = int(Postid)
	db.AddCommatToNotif(commant)
	json.NewEncoder(w).Encode(commant)
}

func (db *Date) AddCommatToNotif(commant utils.Commant) {
	query := `INSERT INTO notif(post_id , user_id , to_user_id , comment_id , action) VALUES (? ,? , ? , ? , ?)
	`
	toUser := db.GetUserIdFromPost(commant.Post_id)
	if toUser == commant.UserId {
		return
	}
	db.DB.Exec(query, commant.Post_id, commant.UserId, toUser, commant.Id, "commeted")
}
