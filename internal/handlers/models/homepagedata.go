package models

type HomePageData struct {
	Posts      []Post
	IsLoggedIn bool
	Username   string
}
