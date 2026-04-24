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

	if err != nil {
		fmt.Errorf("failed to create tables: %v", err)
		return
	} else {
		fmt.Println("post table created successfully")
	}
}
