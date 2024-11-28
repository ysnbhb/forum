package handul

func (db *Date) Insert(user User) error {
	query := `INSERT INTO user (user_name , email , passwd) 
		VALUES (?, ? , ?)
	`
	_, err := db.DB.Exec(query, user.User_name, user.Email, user.Passwd)
	return err
}
