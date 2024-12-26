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
	reaction := utils.Reaction{}
	reaction.Type = r.FormValue("type")
	postid, err := strconv.Atoi(r.FormValue("postid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid post id"})
		return
	}
	like, err := db.CheckLIke(postid, id, "likes", "post_id")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "error in server"})
		return
	}
	dilike, err := db.CheckLIke(postid, id, "dislikes", "post_id")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "error in server"})
		return
	}
	if reaction.Type == "likes" {
		if dilike {
			db.UpdateLike(postid, id, "post_id", "likes")
		} else if like {
			db.DelecttLike(postid, id, "post_id")
		} else {
			db.InsertLike(postid, id, "post_id", "likes")
		}
	} else if reaction.Type == "dislikes" {
		if like {
			db.UpdateLike(postid, id, "post_id", "dislikes")
		} else if dilike {
			db.DelecttLike(postid, id, "post_id")
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
	reaction := utils.Reaction{}
	reaction.Type = r.FormValue("type")
	postid, err := strconv.Atoi(r.FormValue("commateId"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid post id"})
		return
	}
	exist := db.CheckPost(postid)
	if exist != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "post not found"})
		return
	}

	like, err := db.CheckLIke(postid, id, "likes", "comment_id")
	fmt.Println(err)

	dilike, err := db.CheckLIke(postid, id, "dislikes", "comment_id")
	fmt.Println(err)
	if reaction.Type == "likes" {
		if dilike {
			db.UpdateLike(postid, id, "comment_id", "likes")
		} else if like {
			db.DelecttLike(postid, id, "comment_id")
		} else {
			db.InsertLike(postid, id, "comment_id", "likes")
		}
	} else if reaction.Type == "dislikes" {
		if like {
			db.UpdateLike(postid, id, "comment_id", "dislikes")
		} else if dilike {
			db.DelecttLike(postid, id, "comment_id")
		} else {
			db.InsertLike(postid, id, "comment_id", "dislikes")
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid format"})
		return
	}
}

func (db *Date) CheckLIke(id, userId int, typ, column string) (bool, error) {
	var like bool

	// Validate `column` to prevent SQL injection
	if column != "post_id" && column != "comment_id" {
		return false, fmt.Errorf("invalid column name: %s", column)
	}
	query := fmt.Sprintf(`
		SELECT EXISTS (
			SELECT 1 
			FROM reaction
			WHERE %s = ? AND 
			user_id = ? AND 
			type = ?
		)
	`, column)
	err := db.DB.QueryRow(query, id, userId, typ).Scan(&like)
	if err != nil {
		return false, fmt.Errorf("query execution error: %w", err)
	}
	return like, nil
}

func (db *Date) UpdateLike(id, userid int, which, typ string) {
	query := fmt.Sprintf(`
		UPDATE  reaction
		SET type = ?
		WHERE %s = ?
		AND user_id = ?
	`, which)
	db.DB.Exec(query, typ, id, userid)
}

func (db *Date) InsertLike(id, userid int, which, typ string) {
	query := fmt.Sprintf(`
		INSERT INTO reaction(%s ,user_id ,type)
		VALUES( ? , ? , ? )
	`, which)
	db.DB.Exec(query, id, userid, typ)
}

func (db *Date) DelecttLike(id, userid int, which string) {
	query := fmt.Sprintf(`
		DELETE FROM reaction WHERE %s = ? AND user_id = ?
	`, which)
	db.DB.Exec(query, id, userid)
}

func (db *Date) CheckPost(id int) error {
	query := `
		SELECT EXISTS (
			SELECT 1 FORM post WHERE id = ?
		)
	`
	exist := false
	return db.DB.QueryRow(query).Scan(&exist)
}
