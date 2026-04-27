package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"forum/internal/auth"
	"forum/internal/db"
	"forum/internal/handlers"

	// "encoding/json"
	"github.com/gofrs/uuid/v5"

	// "encoding/json"

	// "io"

	_ "github.com/mattn/go-sqlite3"
)

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
	http.ServeFile(w, r, "ui/templates/home.html")
}

func handlePostPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "ui/templates/post.html")
}

func main() {
	database, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		fmt.Errorf("failed to open database %v", err)
	}
	defer database.Close()

	db.Setup(database)

	// 3. Create handler with dependency injection
	a := &auth.AuthHandler{DB: database}
	h := &handlers.Handler{DB: database, Auth: a}

	mux := http.NewServeMux()
	// mux.HandleFunc("/{$}", root)
	mux.HandleFunc("/", h.Home)
	mux.HandleFunc("/registering", handleRegisterHtml)
	mux.HandleFunc("/ping", ping)
	mux.HandleFunc("/register", a.Register)
	mux.HandleFunc("/log", handleLoginPage)
	mux.HandleFunc("/login", a.Login)
	mux.HandleFunc("/comment", h.Comment)
	mux.HandleFunc("/sendpost", h.SendPost)
	mux.HandleFunc("/post", handlePostPage)
	mux.HandleFunc("/like", h.Like)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static/"))))
	fmt.Println("server running on port 8080")

	u4, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("failed to generate unique id %v", err)
	}

	fmt.Println(u4)

	if err != nil {
		fmt.Errorf("failed to create tables: %v", err)
		return
	} else {
		fmt.Println("user created successfully")
	}

	log.Fatal(http.ListenAndServe(":8080", mux))
}
