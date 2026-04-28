package models

type Post struct {
	ID           int
	Title        string
	Content      string
	Username     string 
	Categories   []string
	Comments     []Comment
	CommentCount int
	LikeCount    int
	DislikeCount int
	Initial 	string
	Created_at   string
}
