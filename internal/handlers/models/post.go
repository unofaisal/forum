package models

type Post struct {
	ID           int
	Title        string
	Content      string
	Comments     []Comment
	CommentCount int
	LikeCount    int
	DislikeCount int
}
