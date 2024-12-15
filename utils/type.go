package utils

type User struct {
	User_name string `json:"user_name"`
	Email     string `json:"email"`
	Passwd    string `json:"passwd"`
}

type ErrorData struct {
	Msg1       string
	Msg2       string
	StatusCode int
}


