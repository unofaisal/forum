package handlers

import (
	"fmt"
	"net/http"
)

func (h *Handler) SendPost(w http.ResponseWriter, r *http.Request) {
	postTitle := r.FormValue("postitle")
	postContent := r.FormValue("postContent")
	user_id := 1
	fmt.Println("post sent successfully")
	if postTitle == "" {
		http.Error(w, "post title is required", http.StatusBadRequest)
		return
	}
	if postContent == "" {
		http.Error(w, "post contentn is required", http.StatusBadRequest)
		return
	}

	schema := `INSERT INTO posts (title, content, user_id) VALUES (?, ?, ?)`

	_, err := h.DB.Exec(schema, postTitle, postContent, user_id)

	if err != nil {
		fmt.Printf("failed to add post into the database: %v", err)
	} else {
		// fmt.Fprintf(w, "successfuly added post into the database")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
