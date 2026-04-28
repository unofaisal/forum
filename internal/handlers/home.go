package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"forum/internal/handlers/models"
)

func (h *Handler) getLikes(postID int) (int, int) {
	var likes int
	var dislikes int

	err := h.DB.QueryRow(
		`SELECT COUNT(*) FROM reactions WHERE post_id = ? AND value = 1`,
		postID,
	).Scan(&likes)
	if err != nil {
		likes = 0
	}

	err = h.DB.QueryRow(
		`SELECT COUNT(*) FROM reactions WHERE post_id = ? AND value = -1`,
		postID,
	).Scan(&dislikes)
	if err != nil {
		dislikes = 0
	}

	return likes, dislikes
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	// http.ServeFile(w, r, "ui/templates/home.html")

	fmt.Println("Home hit")
	filterCategory := r.URL.Query().Get("category")

	var row *sql.Rows
	var err error

	// 2. Decide which query to run based on whether a filter exists
	if filterCategory != "" {
		schemaFilterGet := `
			SELECT 
    p.id, 
    p.title, 
    p.content, 
    u.username
FROM posts p
JOIN post_categories pc ON p.id = pc.post_id
JOIN categories c ON pc.category_id = c.id
LEFT JOIN users u ON p.user_id = u.id
WHERE c.name = ?`
		row, err = h.DB.Query(schemaFilterGet, filterCategory)
	} else {
		schemaPostGet := `SELECT 
    p.id, 
    p.title, 
    p.content, 
    u.username
FROM posts p
LEFT JOIN users u ON p.user_id = u.id`
		row, err = h.DB.Query(schemaPostGet)
	}

	if err != nil {
		http.Error(w, "failed to load posts "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer row.Close()

	var post []models.Post

	for row.Next() {
		var p models.Post

		err := row.Scan(&p.ID, &p.Title, &p.Content, &p.Username)
		if err != nil {
			continue
		}

		categoryQuery := `
			SELECT c.name FROM categories c
			JOIN post_categories pc ON c.id = pc.category_id
			WHERE pc.post_id = ?`

		catRows, err := h.DB.Query(categoryQuery, p.ID)
		if err == nil {
			for catRows.Next() {
				var catName string
				if err := catRows.Scan(&catName); err == nil {
					p.Categories = append(p.Categories, catName)
				}
			}
			catRows.Close()
		}

		countCommentQuery := `SELECT COUNT(*) FROM comments WHERE post_id = ?`

		err = h.DB.QueryRow(countCommentQuery, p.ID).Scan(&p.CommentCount)
		if err != nil {
			p.CommentCount = 0
		}
		commentsQuery := `SELECT 
	c.id, 
	c.comment, 
	c.user_id, 
	c.post_id,
	u.username
FROM comments c
LEFT JOIN users u ON c.user_id = u.id
WHERE c.post_id = ?`

		commentRows, err := h.DB.Query(commentsQuery, p.ID)
		if err == nil {
			for commentRows.Next() {
				var c models.Comment
				err := commentRows.Scan(
					&c.ID,
					&c.Comment,
					&c.User_id,
					&c.Post_id,
					&c.Username,
				)
				if err != nil {
					log.Println("scan error:", err)
					continue
				}
				if len(c.Username) > 0 {
					c.Initial = string(c.Username[0])
				} else {
					c.Initial = "?"
				}
				p.Comments = append(p.Comments, c)
			}
			commentRows.Close()
		}

		var Initial string

		if len(p.Username) > 0 {
			Initial = string(p.Username[0])
		} else {
			Initial = "?"
		}
		fmt.Println(Initial)

		likes, dislikes := h.getLikes(p.ID)

		p.Initial = Initial
		p.LikeCount = likes
		p.DislikeCount = dislikes

		post = append(post, p)
		fmt.Println(post)
	}

	fmt.Printf("Rendering %d posts\n", len(post))
	userID, ok := h.Auth.GetUserIDFromSession(r)

	var username string
	if ok {
		err := h.DB.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
		if err != nil {
			username = ""
		}
	}

	data := models.HomePageData{
		Posts:      post,
		IsLoggedIn: ok,
		Username:   username,
	}

	h.RenderTemplate(w, "home", data)
}
