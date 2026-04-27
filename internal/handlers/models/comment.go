package models

type Comment struct {
	ID       int
	Comment  string
	User_id   int
	Post_id  int
	Username string
	Initial  string
}