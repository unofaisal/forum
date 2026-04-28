package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

type mockAuth struct {
	userID int
	ok     bool
}

func (m *mockAuth) GetUserIDFromSession(r *http.Request) (int, bool) {
	return m.userID, m.ok
}

func makeCommentRequest(fields map[string]string) *http.Request {
	form := url.Values{}
	for k, v := range fields {
		form.Set(k, v)
	}
	req := httptest.NewRequest("POST", "/comment", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

func TestComment_Unauthorized(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	h := &Handler{
		DB:   db,
		Auth: &mockAuth{ok: false},
	}

	w := httptest.NewRecorder()
	h.Comment(w, makeCommentRequest(map[string]string{
		"content": "nice post",
		"post_id": "1",
	}))

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "unauthorized") {
		t.Errorf("unexpected body: %s", w.Body.String())
	}
}

func TestComment_EmptyContent(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	h := &Handler{
		DB:   db,
		Auth: &mockAuth{ok: true, userID: 1},
	}

	w := httptest.NewRecorder()
	h.Comment(w, makeCommentRequest(map[string]string{
		"content": "",
		"post_id": "1",
	}))

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "comment cannot be empty") {
		t.Errorf("unexpected body: %s", w.Body.String())
	}
}

func TestComment_InvalidPostID(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	h := &Handler{
		DB:   db,
		Auth: &mockAuth{ok: true, userID: 1},
	}

	w := httptest.NewRecorder()
	h.Comment(w, makeCommentRequest(map[string]string{
		"content": "nice post",
		"post_id": "notanumber",
	}))

	// currently the handler ignores the strconv error and postID becomes 0
	// this test documents the current behaviour — ideally should be 400
	if w.Code == http.StatusUnauthorized {
		t.Errorf("should not be unauthorized for invalid post_id")
	}
}

func TestComment_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	h := &Handler{
		DB:   db,
		Auth: &mockAuth{ok: true, userID: 1},
	}

	mock.ExpectExec("INSERT INTO comments").
		WithArgs("nice post", 1, 1).
		WillReturnError(fmt.Errorf("db error"))

	w := httptest.NewRecorder()
	h.Comment(w, makeCommentRequest(map[string]string{
		"content": "nice post",
		"post_id": "1",
	}))

	// handler only prints the error and doesn't write an HTTP error response
	// so the client gets a blank 200 — documenting current behaviour
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 (no http error written on db fail), got %d", w.Code)
	}
}

func TestComment_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	h := &Handler{
		DB:   db,
		Auth: &mockAuth{ok: true, userID: 1},
	}

	mock.ExpectExec("INSERT INTO comments").
		WithArgs("nice post", 1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	w := httptest.NewRecorder()
	h.Comment(w, makeCommentRequest(map[string]string{
		"content": "nice post",
		"post_id": "1",
	}))

	if w.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", w.Code)
	}
	if w.Header().Get("Location") != "/" {
		t.Errorf("expected redirect to /, got %s", w.Header().Get("Location"))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled mock expectations: %v", err)
	}
}