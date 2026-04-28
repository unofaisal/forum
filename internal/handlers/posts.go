package handlers

import (
	"fmt"
	"net/http"
)

func (h *Handler) SendPost(w http.ResponseWriter, r *http.Request) {
	postTitle := r.FormValue("postitle")
	postContent := r.FormValue("postContent")
	user_id, ok := h.Auth.GetUserIDFromSession(r)
	if !ok {
		http.Redirect(w, r, "/log", http.StatusSeeOther)
		// http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
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

	result, err := h.DB.Exec(schema, postTitle, postContent, user_id)

	if err != nil {
		fmt.Printf("failed to add post into the database: %v", err)
	} else {
		postID, _ := result.LastInsertId()
		r.ParseForm()
		selectedCategories := r.Form["category"]

		for _, catName := range selectedCategories {
			var catID int
			err := h.DB.QueryRow("SELECT id FROM categories WHERE name = ?", catName).Scan(&catID)
			if err == nil {
				h.DB.Exec("INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)", postID, catID)
			}
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
