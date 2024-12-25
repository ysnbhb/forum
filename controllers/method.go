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

func (db *Date) Insert(user utils.User) (int, error) {
	query := `INSERT INTO user (user_name , email , passwd) 
		VALUES (?, ? , ?)
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("sorry but there are error in server try anther time")
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.User_name, user.Email, user.Passwd)
	if err == nil {
		return http.StatusOK, nil
	}
	if strings.Contains(err.Error(), "user_name") {
		return http.StatusFound, fmt.Errorf("user name already used try anther user name")
	} else if strings.Contains(err.Error(), "email") {
		return http.StatusFound, fmt.Errorf("email already used try anther email")
	}
	return http.StatusInternalServerError, fmt.Errorf("sorry but there are error in server try anther time")
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

func (db *Date) CraeteSession(userid int, session string) error {
	query := `INSERT INTO session(user_id , uid)
		VALUES(?,?)
		ON CONFLICT DO UPDATE SET uid = ?
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(userid, session, session)

	return err
}

func (db *Date) TakeName(w http.ResponseWriter, str string) bool {
	id := db.TakeId(str)
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
