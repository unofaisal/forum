package handlers

import (
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
	schemaPostGet := `SELECT id, title, content FROM posts`

	row, err := h.DB.Query(schemaPostGet)
	if err != nil {
		http.Error(w, "failed to load posts", http.StatusInternalServerError)
		return
	}

	defer row.Close()

	var post []models.Post

	for row.Next() {
		var p models.Post

		err := row.Scan(&p.ID, &p.Title, &p.Content)
		if err != nil {
			continue
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

		likes, dislikes := h.getLikes(p.ID)

		p.LikeCount = likes
		p.DislikeCount = dislikes

		post = append(post, p)
		fmt.Println(post)
	}

	fmt.Printf("Rendering %d posts\n", len(post))

	h.RenderTemplate(w, "home", post)
}
