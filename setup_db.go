package main

import (
	"database/sql"
	"fmt"
)

func setup() {
	var err error
	db, err = sql.Open("sqlite3", "forum.db")
	if err != nil {
		fmt.Errorf("failed to open database %v", err)
	}

	schemaComment := `
	CREATE TABLE IF NOT exists comments(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		comment TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		post_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
			)`

	_, err = db.Exec(schemaComment)

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
	_, err = db.Exec(schemaPost)

	schemaLikes := `
	CREATE TABLE IF NOT EXISTS reactions(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		post_id INTEGER NOT NULL,
		value INTEGER NOT NULL, --1 = like, -1 = dislike
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(user_id, post_id)
		)`
	_, err = db.Exec(schemaLikes)

	if err != nil {
		fmt.Errorf("failed to create tables: %v", err)
		return
	} else {
		fmt.Println("posts table created successfully")
	}
}
