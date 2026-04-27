package handlers

import (
	"fmt"
	"net/http"
	"strconv"
)

func (h *Handler) Comment(w http.ResponseWriter, r *http.Request) {
	content := r.FormValue("content")
	postID, err := strconv.Atoi(r.FormValue("post_id"))
	user_id := 1
	// created_at := time.Now()

	if content == "" {
		http.Error(w, "comment cannot be empty", http.StatusBadRequest)
		return
	}
	schemaComment := `INSERT INTO comments (comment, user_id, post_id) VALUES (?, ?, ?)`

	_, err = h.DB.Exec(schemaComment, content, user_id, postID)

	if err != nil {
		fmt.Printf("failed to add comment into the database: %v", err)
		return
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
