package db

import (
	"database/sql"
	"fmt"
)


<<<<<<< HEAD:internal/db/setup_db.go
func Setup(database *sql.DB) {
	
	schema := `
	CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
		`
	_, err := database.Exec(schema)

=======
>>>>>>> a2cf5b4 (modularise code):setup_db.go
	schemaComment := `
	CREATE TABLE IF NOT exists comments(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		comment TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		post_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
			)`

<<<<<<< HEAD:internal/db/setup_db.go
	_, err = database.Exec(schemaComment)
=======
	_, err = db.Exec(schemaComment)
>>>>>>> a2cf5b4 (modularise code):setup_db.go

	if err != nil {
		fmt.Errorf("failed to create tables: %v", err)
		return
	} else {
		fmt.Println("comments table created successfully")
	}

	schemaPost := `
	CREATE TABLE IF NOT EXISTS posts(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		user_id
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
		`
	_, err = database.Exec(schemaPost)

	schemaLikes := `
	CREATE TABLE IF NOT EXISTS reactions(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		post_id INTEGER NOT NULL,
		value INTEGER NOT NULL, --1 = like, -1 = dislike
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(user_id, post_id)
		)`
<<<<<<< HEAD:internal/db/setup_db.go
	_, err = database.Exec(schemaLikes)
=======
	_, err = db.Exec(schemaLikes)
>>>>>>> a2cf5b4 (modularise code):setup_db.go

	if err != nil {
		fmt.Errorf("failed to create tables: %v", err)
		return
	} else {
		fmt.Println("posts table created successfully")
	}
}
