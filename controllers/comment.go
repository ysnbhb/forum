package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"forum/utils"
)

func (db *Date) GetCommat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": http.StatusText(http.StatusMethodNotAllowed)})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	postid := r.URL.Query().Get("postid")
	offset := r.URL.Query().Get("offset")
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
	row, err := (db.DB.Query(query, postid, 5, offset))
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
		comment.Reaction = db.GetReaction(r, comment.Id, "comment_id")
		comments = append(comments, comment)
	}
	json.NewEncoder(w).Encode(comments)
}
