package main

import (
	"bytes"
	"database/sql"

	// "encoding/json"

	// "encoding/json"
	"fmt"
	// "io"
	"log"
	"log/slog"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type UserData struct {
	Name string
}

func root(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "ui/templates/home.html")

	_, err := w.Write([]byte("this is the home page\n"))
	if err != nil {
		slog.Error("error writing the response")
	}

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

	Firstname := r.FormValue("name")
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	// userList := params["name"]

	// fmt.Println(params)
	if Firstname == "" || username == "" || password == "" || email == "" {
		http.Error(w, "all fields are required", http.StatusBadRequest)
		return
	}

	schema := `
	INSERT INTO users (username, password_hash) VALUES (?, ?)`

	user := "user"
	name := "unknown"
	pass := "unknown"
	mail := "unknown"
	// var output bytes.Buffer

	// output.WriteString("Welcome ")

	if len("unknown") > 0 {
		user = username
		pass = password
		mail = email
		name = Firstname
	}

	_, err := db.Exec(schema, username, password)

	if err != nil {
		fmt.Println("DB error: ", err)
		http.Error(w, "failed to save data into databse username already exists", http.StatusInternalServerError)
		return
	} else {
		fmt.Fprintln(w, "data saved successfully into database", http.StatusOK)
	}

	// _, err := w.Write(output.Bytes())

	fmt.Fprintf(w, "welcome %s\nUsername: %s\nEmail: %s\nPassword: %s\n", name, user, mail, pass)
	// body, diode := io.ReadAll(r.Body)

	var data UserData

	// err := json.Unmarshal([]byte(name), &data.Name)
	// diode := json.NewDecoder(r.Body).Decode(&data)

	// if err != nil {
	// 	fmt.Errorf("failed to get userdata %v", err)
	// 	return
	// }

	fmt.Println(name)
	// fmt.Println(&db)
	fmt.Println(data.Name)
}

func login(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("username")
	password := r.FormValue("password")

	if user == "" || password == "" {
		http.Error(w, "all fields must be filled", http.StatusBadRequest)
		return
	}

	schema := `
	SELECT username, password_hash FROM users WHERE username = ?
	`

	row := db.QueryRow(schema, user)

	// if row == "" {
	// 	http.Error(w, "user not found", http.StatusNotFound)
	// 	return
	// }

	// if err != nil {
	// 	http.Error(w, "user not found", http.StatusNotFound)
	// 	return
	// }

	var dbusername, dbpassword string

	row.Scan(&dbusername, &dbpassword)
	if dbpassword != password {
		http.Error(w, "user unknown try again", http.StatusForbidden)
		return
	}
	fmt.Println(dbusername, user)

	// for row.Next() {
	// 	// var username, password string
	// 	row.Scan(&username, &password)

	fmt.Fprintf(w, "Welcome back %v", dbusername)
	// }
}

func handleLoginPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "ui/templates/login.html")
}

func handleRegisterHtml(w http.ResponseWriter, r *http.Request) {
	// fs := http.FileServer(http.Dir("./ui/templates/register.html"))
	http.ServeFile(w, r, "ui/templates/register.html")
}

func handleHomePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/ui/templates/home.html")
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	schema := `
	SELECT username, password_hash FROM users`
	row, err := db.Query(schema)
	if err != nil {
		http.Error(w, "failed to retrieve data from the database", http.StatusInternalServerError)
		return
	}
	// defer db.Close()

	for row.Next() {
		var username, password string
		row.Scan(&username, &password)

		fmt.Fprintf(w, "username: %v, password: %v\n", username, password)
	}
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
	mux.HandleFunc("/login", login)

	fmt.Println("server running on port 8080")
	var err interface{}

	db, err = sql.Open("sqlite3", "forum.db")
	if err != nil {
		fmt.Errorf("failed to open database %v", err)
	}

	schema := `
	CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
		`
	_, err = db.Exec(schema)

	if err != nil {
		fmt.Errorf("failed to create tables: %v", err)
		return
	} else {
		fmt.Println("user created successfully")
	}

	defer db.Close()
	log.Fatal(http.ListenAndServe(":8080", mux))
}
