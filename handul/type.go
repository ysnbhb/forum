package handul

import "database/sql"

type Date struct {
	DB *sql.DB
}

type User struct {
	User_name string `json:"user_name"`
	Email     string `json:"email"`
	Passwd    string `json:"passwd"`
}
