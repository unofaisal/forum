package main

import (
	"bytes"
	"database/sql"
	"strings"

	// "crypto/des"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoot(t *testing.T) {
	// req := httptest(http.MethodGet, "/root", nil)
	w := httptest.NewRecorder()
	root(w, nil)

	desiredcode := http.StatusOK

	if w.Code != http.StatusOK {
		t.Errorf("bad request expected: %v but returned: %v", desiredcode, w.Code)
	}

	expectedMessage := "this is the home page\n"

	if !bytes.Equal([]byte(expectedMessage), w.Body.Bytes()) {
		t.Errorf("bad response expected: %q but got: %q", expectedMessage, w.Body.String())
	}
}

func TestPing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ping?user=jack", nil)

	w := httptest.NewRecorder()

	ping(w, req)

	desiredcode := http.StatusOK

	if desiredcode != w.Code {
		t.Errorf("Bad response expected: %v but received: %v", desiredcode, w.Code)
	}
	expectedMessage := []byte("Hello jack!\n")

	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("Bad response body expected: %q but received: %v", expectedMessage, w.Body.String())
	}
}

func TestRegister(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/register", nil)

	// Firstname:= "name"
	// username:= "username"
	// email := "email"
	// password:= "password"

	desiredcode := http.StatusOK

	w := httptest.NewRecorder()
	register(w, req)
	if w.Code != desiredcode {
		t.Errorf("bad response expected: %v but received: %v body %v", desiredcode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("Welcome \nusername: \n email: \n password: !\n")

	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad body expected: %q but receives: %q", expectedMessage, w.Body.String())
	}
}

func TestGetUser(t *testing.T) {
	db, _ = sql.Open("sqlite3", ":memory:")
	schema := `
	 CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	 `
	db.Exec(schema)

	_, _ = db.Exec(
		"INSERT INTO users (username, password_hash) VALUES (?, ?)",
		"testuser",
		"123",
	)
	reg := httptest.NewRequest(http.MethodGet, "/register", nil)
	req := httptest.NewRequest(http.MethodGet, "/getusers", nil)
	w1 := httptest.NewRecorder()
	register(w1, reg)

	w2 := httptest.NewRecorder()
	getUsers(w2, req)

	register(w1, reg)
	getUsers(w2, req)
	// if err != nil {
	// 	t.Error("failed to create fake data")
	// }

	desiredcode := http.StatusOK

	if w2.Code != desiredcode {
		t.Errorf("Error bad request expected %v but received: %v", desiredcode, w2.Code)
	}

	body := w2.Body.String()

	if !strings.Contains(body, "testuser") {
		t.Errorf("expected response to contain testuser, got: %s", body)
	}
}
