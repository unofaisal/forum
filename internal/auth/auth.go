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
	// http.ServeFile(w, r, "ui/templates/register.html")

	// Firstname := r.FormValue("name")
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confPassword := r.FormValue("confirmPassword")

	// userList := params["name"]

	// fmt.Println(params)
	if confPassword == "" || username == "" || password == "" || email == "" {
		http.Error(w, "all fields are required", http.StatusBadRequest)
		return
	}

	if confPassword != password {
		http.Error(w, "passwords do not match", http.StatusBadRequest)
		// confPassError = "passwords do not match"
		return
	}
	// confPassError = r.FormValue("confirmPasswordError")

	schema := `
	INSERT INTO users (username, password_hash, email) VALUES (?, ?, ?)`

	
	// name = Firstname
	// }

	passByte := []byte(password)

	fmt.Println(passByte)

	hashedPassword, error := bcrypt.GenerateFromPassword(passByte, bcrypt.DefaultCost)

	if error != nil {
		http.Error(w, "failed to hash the password", http.StatusInternalServerError)
		return
		// panic(error)
	}
	
result, err := h.DB.Exec(schema, username, string(hashedPassword), email)
if err != nil {
    http.Error(w, "failed to save data email or username is already taken", http.StatusInternalServerError)
    return
}

userID, err := result.LastInsertId()
if err != nil {
    http.Error(w, "failed to get user id", http.StatusInternalServerError)
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
`, sessionID, userID)

if err != nil {
    http.Error(w, "failed to create session", http.StatusInternalServerError)
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



	http.Redirect(w, r, "/", http.StatusSeeOther)
	return

	// _, err := w.Write(output.Bytes())

	// body, diode := io.ReadAll(r.Body)

	// var data UserData

	// err := json.Unmarshal([]byte(name), &data.Name)
	// diode := json.NewDecoder(r.Body).Decode(&data)

	// if err != nil {
	// 	fmt.Errorf("failed to get userdata %v", err)
	// 	return
	// }

	// fmt.Println(name)
	// fmt.Println(&db)
	// fmt.Println(data.Name)
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


	// for row.Next() {
	// 	// var username, password string
	// 	row.Scan(&username, &password)

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
