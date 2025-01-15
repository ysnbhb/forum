package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"forum/utils"

	"golang.org/x/crypto/bcrypt"
)

func (db *Date) Insert(user utils.User, typeOflog string) (int, int64, error) {
	query := `INSERT INTO user (user_name , email , passwd , typeOflog) 
		VALUES (?, ? , ?)
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return http.StatusInternalServerError, 0, fmt.Errorf("sorry but there are error in server try anther time")
	}
	defer stmt.Close()
	res, err := stmt.Exec(user.User_name, user.Email, user.Passwd, typeOflog)
	if err == nil {
		lastId, _ := res.LastInsertId()
		return http.StatusOK, lastId, nil
	}
	if strings.Contains(err.Error(), "user_name") {
		return http.StatusFound, 0, fmt.Errorf("user name already used try anther user name")
	} else if strings.Contains(err.Error(), "email") {
		return http.StatusFound, 0, fmt.Errorf("email already used try anther email")
	}
	return http.StatusInternalServerError, 0, fmt.Errorf("sorry but there are error in server try anther time")
}

func (db *Date) Select(userIfo, passwd string) (int, error) {
	query := `SELECT id , passwd FROM user
		WHERE user_name = ? OR email = ?
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
		return -2, fmt.Errorf("problem in server try anther time")
	}
	defer stmt.Close()
	var hashpasswd string
	var id int
	err = stmt.QueryRow(userIfo, userIfo).Scan(&id, &hashpasswd)
	if err == sql.ErrNoRows {
		return -1, fmt.Errorf("user or password not correct")
	} else if err != nil {
		return -1, fmt.Errorf("check your input")
	}
	if bcrypt.CompareHashAndPassword([]byte(hashpasswd), []byte(passwd)) != nil {
		return -1, fmt.Errorf("user or password not correct")
	}
	return id, nil
}

func (db *Date) SelectAuth(userIfo string) (int, string) {
	query := `SELECT id , typeOfLog FROM user
		WHERE  email = ?
	`
	var typeOfLog string
	var id int
	err := db.DB.QueryRow(query, userIfo, userIfo).Scan(&id, &typeOfLog)
	if err == sql.ErrNoRows {
		return -1, ""
	} else if err != nil {
		return -1, ""
	}
	return id, typeOfLog
}

func (db *Date) CraeteSession(userid int, session string) error {
	query := `INSERT INTO session(user_id , uid)
	VALUES(?,?)
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(userid, session)
	if err != nil {
		cookie := ""
		query = `SELECT uid FROM session WHERE user_id = ?`
		db.DB.QueryRow(query, userid).Scan(&cookie)
		query = `UPDATE session SET uid = ? , create_date = CURRENT_TIMESTAMP WHERE user_id = ?
		`
		stmt, err := db.DB.Prepare(query)
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(session, userid)

		delete(utils.RateLimitData, cookie)
		return err
	}
	return nil
}

func (db *Date) TakeName(w http.ResponseWriter, r *http.Request) bool {
	id := db.GetIdFromReq(r)
	userName, err := db.GetName(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "user not Found"})
		return false
	}
	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(map[string]string{"userName": userName})
	return true
}

func (db *Date) GetName(id int) (string, error) {
	userName := ""
	query := `SELECT (user_name) FROM user WHERE id = ?`
	err := db.DB.QueryRow(query, id).Scan(&userName)
	return userName, err
}

func (db *Date) LastID(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT seq FROM sqlite_sequence WHERE name = 'post';

	`
	var lastId int
	db.DB.QueryRow(query).Scan(&lastId)
	fmt.Fprint(w, lastId)
}

func (db *Date) TakeId(secion string) int {
	query := `
		SELECT user_id FROM session WHERE uid = ?
	`
	id := 0
	err := db.DB.QueryRow(query, secion).Scan(&id)
	if err != nil {
		return -1
	}
	return id
}

func (db *Date) DelectSeoin(token string) error {
	query := `
		DELETE FROM session WHERE user_id = ?
	`
	_, err := db.DB.Exec(query, token)
	return err
}

func (db *Date) DelectNotifLIke(post_id, user_id int) {
	query := `DELETE FROM notif  WHERE post_id = ? AND user_id = ? AND reaction_id IS NOT NULL`
	db.DB.Exec(query, post_id, user_id)
}

func (db *Date) UpdateNotifLIke(typ string, post_id, user_id int) {
	query := `UPDATE notif SET action = ? WHERE post_id = ? AND user_id = ? AND reaction_id IS NOT NULL`
	_, err := db.DB.Exec(query, typ, post_id, user_id)
	fmt.Println(err)
}

func (db *Date) InsertNotifLIke(typ string, post_id, user_id, to_user_id, notif_id int) {
	if user_id == to_user_id {
		return
	}
	query := `INSERT INTO notif(post_id, user_id , to_user_id, action , reaction_id) 
		VALUES (? , ? , ? , ? , ?)`
	db.DB.Exec(query, post_id, user_id, to_user_id, typ, notif_id)
}

func (db *Date) GetIdFromReq(r *http.Request) int {
	cookie, err := r.Cookie("token")
	if err != nil {
		return -1
	}
	userId := db.TakeId(cookie.Value)
	return userId
}

func (db *Date) TakeIDfromAuth(str string) bool {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM noLOg WHERE token  = ? 
		)
	`
	id := false
	err := db.DB.QueryRow(query, str).Scan(&id)
	if err != nil {
		return false
	}
	return id
}
