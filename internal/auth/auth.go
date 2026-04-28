package auth

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB *sql.DB
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confPassword := r.FormValue("confirmPassword")

	if confPassword == "" || username == "" || password == "" || email == "" {
		http.Error(w, "all fields are required", http.StatusBadRequest)
		return
	}

	if confPassword != password {
		http.Error(w, "passwords do not match", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed to hash the password", http.StatusInternalServerError)
		return
	}

	_, err = h.DB.Exec(
		`INSERT INTO users (username, password_hash, email) VALUES (?, ?, ?)`,
		username, string(hashedPassword), email,
	)
	if err != nil {
		fmt.Println("DB error: ", err)
		http.Error(w, "failed to save data into database, username or email already exists", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		http.Error(w, "all fields must be filled", http.StatusBadRequest)
		return
	}

	schema := `
	SELECT username, id, password_hash FROM users WHERE email = ?
	`

	row := h.DB.QueryRow(schema, email)

	// if row == "" {
	// 	http.Error(w, "user not found", http.StatusNotFound)
	// 	return
	// }

	// if err != nil {
	// 	http.Error(w, "user not found", http.StatusNotFound)
	// 	return
	// }
	var dbemail, dbpassword string
	var user_id int

	row.Scan(&dbemail, &user_id, &dbpassword)

	err := bcrypt.CompareHashAndPassword([]byte(dbpassword), []byte(password))
	if err != nil {
		http.Error(w, "user unknown try again", http.StatusForbidden)
		return
	}
	u4, err := uuid.NewV4()
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	sessionID := u4.String()

	_, err = h.DB.Exec(`
		INSERT INTO sessions (sessionId, user_id, expires_at)
		VALUES (?, ?, datetime('now', '+1 hour'))
	`, sessionID, user_id)
	if err != nil {
		fmt.Println("DB error: ", err)
		http.Error(w, "failed to save session", http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:     "session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   3600,
	}
	http.SetCookie(w, cookie)

	fmt.Fprintf(w, "Welcome back %v", dbemail)
	fmt.Println(dbpassword)

	// if dbpassword != password {
	// 	http.Error(w, "user unknown try again", http.StatusForbidden)
	// 	return
	// }

	fmt.Println(dbemail, user)

	// for row.Next() {
	// 	// var username, password string
	// 	row.Scan(&username, &password)

	fmt.Fprintf(w, "Welcome back %v", dbemail)
	// }
}

func (h *AuthHandler) GetUserIDFromSession(r *http.Request) (int, bool) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return 0, false
	}

	var userID int
	var expires string

	err = h.DB.QueryRow(`
		SELECT user_id, expires_at 
		FROM sessions 
		WHERE sessionId = ?
	`, cookie.Value).Scan(&userID, &expires)
	if err != nil {
		fmt.Println("error fetching user id from session: ", err.Error())
		return 0, false
	}
	fmt.Println("user id from session: ", userID)

	return userID, true
}
