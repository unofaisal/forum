package models

type UserData struct {
	Name string
}

type Comment struct {
	ID       int
	Comment  string
	User_id  int
	Post_id  int
	Username string
	Initial string
}

type Post struct {
	ID           int
	Title        string
	Content      string
	Comments     []Comment
	CommentCount int
	LikeCount    int
	DislikeCount int
}