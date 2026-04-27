package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"

	// "encoding/json"
	"github.com/gofrs/uuid/v5"

	// "encoding/json"
	"golang.org/x/crypto/bcrypt"
	// "io"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type UserData struct {
	Name string
}

type Comment struct {
	ID       int
	Comment  string
	User_id  int
	Post_id  int
	Username string
	Initial string
}

type Post struct {
	ID           int
	Title        string
	Content      string
	Comments     []Comment
	CommentCount int
	LikeCount    int
	DislikeCount int
}

func root(w http.ResponseWriter, r *http.Request) {
	schemaPostGet := `SELECT id, title, content FROM posts`

	row, err := db.Query(schemaPostGet)
	if err != nil {
		http.Error(w, "failed to load posts", http.StatusInternalServerError)
		return
	}

	defer row.Close()

	var post []Post

	for row.Next() {
		var p Post

		err := row.Scan(&p.ID, &p.Title, &p.Content)
		if err != nil {
			continue
		}
		countCommentQuery := `SELECT COUNT(*) FROM comments WHERE post_id = ?`

		err = db.QueryRow(countCommentQuery, p.ID).Scan(&p.CommentCount)
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

		commentRows, err := db.Query(commentsQuery, p.ID)
		if err == nil {
			for commentRows.Next() {
				var c Comment
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
				}else {
					c.Initial = "?"
				}
				p.Comments = append(p.Comments, c)
			}
			commentRows.Close()
		}

		likes, dislikes := getLikes(p.ID)

		p.LikeCount = likes
		p.DislikeCount = dislikes

		post = append(post, p)
		fmt.Println(post)
	}

	tmpl, err := template.ParseFiles("ui/templates/home.html")
	if err != nil {
		fmt.Println("post error: %v", err)
		http.Error(w, "failed to update ui %v", http.StatusNotFound)
		return
	}
	tmpl.Execute(w, post)
}

func ping(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	userList := params["user"]

	var output bytes.Buffer

	user := "user"
	output.WriteString("Hello ")
	if len(userList) > 0 {
		user = userList[0]
	}
	output.WriteString(user)
	output.WriteString("!\n")

	_, err := w.Write(output.Bytes())
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	// http.ServeFile(w, r, "ui/templates/register.html")

	// Firstname := r.FormValue("name")
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confPassword := r.FormValue("confirmPassword")
	confPassError := r.FormValue("confirmPasswordError")

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

	user := "user"
	// name := "unknown"
	pass := "unknown"
	mail := "unknown"
	// var output bytes.Buffer

	// output.WriteString("Welcome ")

	// if len("unknown") > 0 {
	user = username
	pass = password
	mail = email
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

	_, err := db.Exec(schema, username, string(hashedPassword), email)

	if err != nil {
		fmt.Println("DB error: ", err)
		http.Error(w, "failed to save data into databse username or email already exists", http.StatusInternalServerError)
		return
	} else {
		fmt.Fprintln(w, "data saved successfully into database", http.StatusOK)
	}

	// _, err := w.Write(output.Bytes())
	fmt.Println("this is the errror: ", confPassError)

	fmt.Fprintf(w, "Username: %s\nEmail: %s\nPassword: %s\n", user, mail, pass)
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

func login(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		http.Error(w, "all fields must be filled", http.StatusBadRequest)
		return
	}

	schema := `
	SELECT username, password_hash FROM users WHERE email = ?
	`

	row := db.QueryRow(schema, email)

	// if row == "" {
	// 	http.Error(w, "user not found", http.StatusNotFound)
	// 	return
	// }

	// if err != nil {
	// 	http.Error(w, "user not found", http.StatusNotFound)
	// 	return
	// }
	var dbemail, dbpassword string

	row.Scan(&dbemail, &dbpassword)

	err := bcrypt.CompareHashAndPassword([]byte(dbpassword), []byte(password))
	if err != nil {
		http.Error(w, "user unknown try again", http.StatusForbidden)
		return
	}
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

func handleLoginPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "ui/templates/login.html")
}

func handleRegisterHtml(w http.ResponseWriter, r *http.Request) {
	// fs := http.FileServer(http.Dir("./ui/templates/register.html"))
	http.ServeFile(w, r, "ui/templates/signup.html")
	// http.Handle("registering", fs)
	// fs.ServeHTTP(w,r)
}

func handleHomePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/ui/templates/home.html")
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	schema := `
	SELECT username, password_hash, email FROM users`
	row, err := db.Query(schema)
	if err != nil {
		http.Error(w, "failed to retrieve data from the database", http.StatusInternalServerError)
		return
	}
	// defer db.Close()

	for row.Next() {
		var username, password, email string
		row.Scan(&username, &password, &email)

		fmt.Fprintf(w, "username: %v, password: %v, email: %v\n", username, password, email)
	}
}

func handlePostPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "ui/templates/post.html")
}

func sendPost(w http.ResponseWriter, r *http.Request) {
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

	_, err := db.Exec(schema, postTitle, postContent, user_id)

	if err != nil {
		fmt.Printf("failed to add post into the database: %v", err)
	} else {
		// fmt.Fprintf(w, "successfuly added post into the database")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func comment(w http.ResponseWriter, r *http.Request) {
	content := r.FormValue("content")
	postID, err := strconv.Atoi(r.FormValue("post_id"))
	user_id := 1
	// created_at := time.Now()

	if content == "" {
		http.Error(w, "comment cannot be empty", http.StatusBadRequest)
		return
	}
	schemaComment := `INSERT INTO comments (comment, user_id, post_id) VALUES (?, ?, ?)`

	_, err = db.Exec(schemaComment, content, user_id, postID)

	if err != nil {
		fmt.Printf("failed to add comment into the database: %v", err)
		return
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func getComments(w http.ResponseWriter, r *http.Request) {
	schema := `
	SELECT content, post_id, user_id FROM comments`
	row, err := db.Query(schema)
	if err != nil {
		http.Error(w, "failed to retrieve data from the database", http.StatusInternalServerError)
		return
	}
	// defer db.Close()

	for row.Next() {
		var content, post_id, user_id string
		row.Scan(&content, &post_id, &user_id)

		fmt.Fprintf(w, "content: %v, post_id: %v, user_id: %v\n", content, post_id, user_id)
	}
}

func getLikes(postID int) (int, int) {
	var likes int
	var dislikes int

	err := db.QueryRow(
		`SELECT COUNT(*) FROM reactions WHERE post_id = ? AND value = 1`,
		postID,
	).Scan(&likes)
	if err != nil {
		likes = 0
	}

	err = db.QueryRow(
		`SELECT COUNT(*) FROM reactions WHERE post_id = ? AND value = -1`,
		postID,
	).Scan(&dislikes)
	if err != nil {
		dislikes = 0
	}

	return likes, dislikes
}

func like(w http.ResponseWriter, r *http.Request) {
	postID, _ := strconv.Atoi(r.FormValue("post_id"))
	value, _ := strconv.Atoi(r.FormValue("value"))

	userID := 1
	var existing int
	err := db.QueryRow(
		`SELECT value FROM reactions WHERE user_id = ? AND post_id = ?`,
		userID, postID,
	).Scan(&existing)

	switch {
	case err == sql.ErrNoRows:
		_, err = db.Exec(
			`INSERT INTO reactions (user_id, post_id, value) VALUES (?, ?, ?)`,
			userID, postID, value,
		)

	case err != nil:
		log.Println("query error:", err)

	default:
		// row exists

		if existing == value {
			_, err = db.Exec(
				`DELETE FROM reactions WHERE user_id = ? AND post_id = ?`,
				userID, postID,
			)
		} else {
			// switch like ↔ dislike
			_, err = db.Exec(
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

func main() {
	mux := http.NewServeMux()
	// mux.HandleFunc("/{$}", root)
	mux.HandleFunc("/", root)
	mux.HandleFunc("/registering", handleRegisterHtml)
	mux.HandleFunc("/ping", ping)
	mux.HandleFunc("/register", register)
	mux.HandleFunc("/log", handleLoginPage)
	mux.HandleFunc("/getusers", getUsers)
	mux.HandleFunc("/getcomments", getComments)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/comment", comment)
	mux.HandleFunc("/sendpost", sendPost)
	mux.HandleFunc("/post", handlePostPage)
	mux.HandleFunc("/like", like)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static/"))))
	fmt.Println("server running on port 8080")
	var err interface{}

	db, err = sql.Open("sqlite3", "forum.db")
	if err != nil {
		fmt.Errorf("failed to open database %v", err)
	}

	u4, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("failed to generate unique id %v", err)
	}

	fmt.Println(u4)

	schema := `
	CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
		`
	_, err = db.Exec(schema)

	setup()

	if err != nil {
		fmt.Errorf("failed to create tables: %v", err)
		return
	} else {
		fmt.Println("user created successfully")
	}

	defer db.Close()
	log.Fatal(http.ListenAndServe(":8080", mux))
}
