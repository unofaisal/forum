package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func makeSendPostRequest(fields map[string]string, categories []string) *http.Request {
	form := url.Values{}
	for k, v := range fields {
		form.Set(k, v)
	}
	for _, cat := range categories {
		form.Add("category", cat)
	}
	req := httptest.NewRequest("POST", "/post", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

func TestSendPost_Unauthorized(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	h := &Handler{
		DB:   db,
		Auth: &mockAuth{ok: false},
	}

	w := httptest.NewRecorder()
	h.SendPost(w, makeSendPostRequest(map[string]string{
		"postitle":    "My Post",
		"postContent": "Some content",
	}, nil))

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "unauthorized") {
		t.Errorf("unexpected body: %s", w.Body.String())
	}
}

func TestSendPost_MissingTitle(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	h := &Handler{
		DB:   db,
		Auth: &mockAuth{ok: true, userID: 1},
	}

	w := httptest.NewRecorder()
	h.SendPost(w, makeSendPostRequest(map[string]string{
		"postContent": "Some content",
	}, nil))

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "post title is required") {
		t.Errorf("unexpected body: %s", w.Body.String())
	}
}

func TestSendPost_MissingContent(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	h := &Handler{
		DB:   db,
		Auth: &mockAuth{ok: true, userID: 1},
	}

	w := httptest.NewRecorder()
	h.SendPost(w, makeSendPostRequest(map[string]string{
		"postitle": "My Post",
	}, nil))

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "post contentn is required") {
		t.Errorf("unexpected body: %s", w.Body.String())
	}
}
