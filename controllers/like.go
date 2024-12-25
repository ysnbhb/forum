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
	like := db.CheckLIke(postid, id, "post_id", "likes")
	dilike := db.CheckLIke(postid, id, "post_id", "dislikes")
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
	like := db.CheckLIke(postid, id, "comment_id", "likes")
	dilike := db.CheckLIke(postid, id, "comment_id", "dislikes")
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

func (db *Date) CheckLIke(id, uerId int, tpy, which string) bool {
	var like bool
	query := fmt.Sprintf(`
		SELECT EXISTS (
			SELECT 1 
			FROM reaction
			WHERE %s = ? AND 
			user_id = ? 
			WHERE type = ?
		)
	`, which)
	db.DB.QueryRow(query, id, uerId, tpy).Scan(&like)
	return like
}

func (db *Date) UpdateLike(id, userid int, which, typ string) {
	query := fmt.Sprintf(`
		UPDATE INTO reaction
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
