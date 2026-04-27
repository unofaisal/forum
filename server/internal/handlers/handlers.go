package handlers

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"golang.org/x/crypto/bcrypt"
)

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
				} else {
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
