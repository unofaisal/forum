package db

import (
	"database/sql"
	"fmt"
)


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

	schemaComment := `
	CREATE TABLE IF NOT exists comments(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		comment TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		post_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
			)`

	_, err = database.Exec(schemaComment)

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
	_, err = database.Exec(schemaLikes)

	if err != nil {
		fmt.Errorf("failed to create tables: %v", err)
		return
	} else {
		fmt.Println("posts table created successfully")
	}
}
