package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"forum/utils"
)

func (db *Date) LikePost(w http.ResponseWriter, r *http.Request) {
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
	reaction := utils.Reaction{}
	reaction.Type = r.FormValue("type")
	postid, err := strconv.Atoi(r.FormValue("postid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid post id"})
		return
	}
	exist := db.CheckPost(id)
	if exist != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid post id"})
		return
	}
	like := db.CheckLIke(postid, id, "likes", "post_id")
	if reaction.Type == "likes" {
		if like == "dislikes" {
			db.UpdateLike(postid, id, "post_id", "likes")
		} else if like == "likes" {
			db.DelecttLike(postid, id, "post_id", like)
		} else {
			db.InsertLike(postid, id, "post_id", "likes")
		}
	} else if reaction.Type == "dislikes" {
		if like == "likes" {
			db.UpdateLike(postid, id, "post_id", "dislikes")
		} else if like == "dislikes" {
			db.DelecttLike(postid, id, "post_id", like)
		} else {
			db.InsertLike(postid, id, "post_id", "dislikes")
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid format"})
		return
	}
}

func (db *Date) LikeCommat(w http.ResponseWriter, r *http.Request) {
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
	reaction := utils.Reaction{}
	reaction.Type = r.FormValue("type")
	postid, err := strconv.Atoi(r.FormValue("commateId"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid commate id"})
		return
	}
	exist := db.CheckCommat(postid)
	if exist != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "commat not found"})
		return
	}
	like := db.CheckLIke(postid, id, "likes", "comment_id")
	if reaction.Type == "likes" {
		if like == "dislikes" {
			db.UpdateLike(postid, id, "comment_id", "likes")
		} else if like == "likes" {
			db.DelecttLike(postid, id, "comment_id", like)
		} else {
			db.InsertLike(postid, id, "comment_id", "likes")
		}
	} else if reaction.Type == "dislikes" {
		if like == "likes" {
			db.UpdateLike(postid, id, "comment_id", "dislikes")
		} else if like == "dislikes" {
			db.DelecttLike(postid, id, "comment_id", like)
		} else {
			db.InsertLike(postid, id, "comment_id", "dislikes")
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid format"})
		return
	}
}

func (db *Date) CheckLIke(id, userId int, typ, column string) string {
	action := ""
	query := fmt.Sprintf(`
		SELECT type FROM reaction
			WHERE %s = ? AND 
			user_id = ?
	`, column)
	db.DB.QueryRow(query, id, userId).Scan(&action)
	return action
}

func (db *Date) UpdateLike(id, userid int, colom, typ string) {
	query := fmt.Sprintf(`
		UPDATE  reaction
		SET type = ?
		WHERE %s = ?
		AND user_id = ?
	`, colom)
	db.DB.Exec(query, typ, id, userid)
	if colom == "post_id" {
		db.UpdateNotifLIke(typ, id, userid)
	}
}

func (db *Date) InsertLike(id, userid int, colom, typ string) {
	query := fmt.Sprintf(`
		INSERT INTO reaction(%s ,user_id ,type)
		VALUES( ? , ? , ? )
	`, colom)
	res, _ := db.DB.Exec(query, id, userid, typ)

	if colom == "post_id" {
		last_id, _ := res.LastInsertId()
		to_userid := db.GetUserIdFromPost(id)
		db.InsertNotifLIke(typ, id, userid, to_userid, int(last_id))
	}
}

func (db *Date) DelecttLike(id, userid int, colom string, action string) {
	query := fmt.Sprintf(`
		DELETE FROM reaction WHERE %s = ? AND user_id = ?
	`, colom)
	db.DB.Exec(query, id, userid)
	if colom == "post_id" {
		db.DelectNotifLIke(id, userid)
	}
}

func (db *Date) CheckPost(id int) error {
	query := `
		SELECT EXISTS (
    SELECT 1 FROM post WHERE id = ? )
	`
	exist := false
	return db.DB.QueryRow(query, id).Scan(&exist)
}

func (db *Date) CheckCommat(id int) error {
	query := `
		SELECT EXISTS (
    SELECT 1 FROM comment WHERE id = ? )
	`
	exist := false
	return db.DB.QueryRow(query, id).Scan(&exist)
}

func (db *Date) GetUserIdFromPost(postId int) int {
	userid := 0
	query := `
    SELECT user_id FROM post WHERE id = ? 
	`
	db.DB.QueryRow(query, postId).Scan(&userid)
	return userid
}
