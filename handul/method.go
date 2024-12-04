package handul

import (
	"fmt"
	"strings"
)

func (db *Date) Insert(user User) error {
	query := `INSERT INTO user (user_name , email , passwd) 
		VALUES (?, ? , ?)
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("check your Input")
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.User_name, user.Email, user.Passwd)
	if err == nil {
		return nil
	}
	if strings.Contains(err.Error(), "user_name") {
		return fmt.Errorf("user name")
	} else if strings.Contains(err.Error(), "email") {
		return fmt.Errorf("email")
	}
	return err
}
