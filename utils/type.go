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

type Post struct {
	Id            int      `json:"id"`
	UserName      string   `json:"userName"`
	Title         string   `json:"title"`
	Contant       string   `json:"contant"`
	ImgUrl        any      `json:"img"`
	Categories    []string `json:"categories"`
	Reaction      Reaction `json:"reaction"`
	Idscategories []int    `json:"-"`
	Date          string   `json:"date"`
}

type Reaction struct {
	PostId     int    `json:"-"`
	CommentId  int    `json:"-"`
	NumLike    int    `json:"numLike"`
	NumDisLike int    `json:"numDisLike"`
	Type       string `json:"type"` // here i can know if user like post or not
}

type Commant struct {
	Id       int      `json:"id"`
	UserName string   `json:"userName"`
	Contant  string   `json:"contant"`
	Reaction Reaction `json:"reaction"`
	Date     string   `json:"date"`
}
