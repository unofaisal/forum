package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"forum/internal/auth"
)

func (h *Handler) Like(w http.ResponseWriter, r *http.Request) {
	postID, _ := strconv.Atoi(r.FormValue("post_id"))
	value, _ := strconv.Atoi(r.FormValue("value"))

	userID, ok := auth.GetUserIDFromSession(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var existing int
	err := h.DB.QueryRow(
		`SELECT value FROM reactions WHERE user_id = ? AND post_id = ?`,
		userID, postID,
	).Scan(&existing)

	switch {
	case err == sql.ErrNoRows:
		_, err = h.DB.Exec(
			`INSERT INTO reactions (user_id, post_id, value) VALUES (?, ?, ?)`,
			userID, postID, value,
		)

	case err != nil:
		log.Println("query error:", err)

	default:
		// row exists

		if existing == value {
			_, err = h.DB.Exec(
				`DELETE FROM reactions WHERE user_id = ? AND post_id = ?`,
				userID, postID,
			)
		} else {
			// switch like ↔ dislike
			_, err = h.DB.Exec(
				`UPDATE reactions SET value = ? WHERE user_id = ? AND post_id = ?`,
				value, userID, postID,
			)
		}
	}

	if err != nil {
		log.Println("reaction write error:", err)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
