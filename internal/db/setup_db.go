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

	schemaSessions := `
	CREATE TABLE IF NOT EXISTS sessions(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		sessionId TEXT NOT NULL UNIQUE,
		expires_at DATETIME NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
	)`
	_, err = database.Exec(schemaSessions)
	if err == nil {
		fmt.Println("sessions table created successfully")
	}

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
		user_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
		`
	_, err = database.Exec(schemaPost)

	schemaCategory := `
	CREATE TABLE IF NOT EXISTS categories(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	)`
	database.Exec(schemaCategory)

	schemaPostCategories := `
	CREATE TABLE IF NOT EXISTS post_categories(
		post_id INTEGER,
		category_id INTEGER,
		PRIMARY KEY (post_id, category_id),
		FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE,
		FOREIGN KEY(category_id) REFERENCES categories(id) ON DELETE CASCADE
	)`
	database.Exec(schemaPostCategories)

	categories := []string{"Technology", "Gaming", "Science", "Music"}
	for _, cat := range categories {
		database.Exec("INSERT OR IGNORE INTO categories (name) VALUES (?)", cat)
	}
	fmt.Println("categories tables created and seeded successfully")

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
		fmt.Printf("failed to create tables: %v\n", err)
		return
	}

}
