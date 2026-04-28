package handlers

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func makeLikeRequest(fields map[string]string) *http.Request {
	form := url.Values{}
	for k, v := range fields {
		form.Set(k, v)
	}
	req := httptest.NewRequest("POST", "/like", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

func TestLike_Unauthorized(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	h := &Handler{
		DB:   db,
		Auth: &mockAuth{ok: false},
	}

	w := httptest.NewRecorder()
	h.Like(w, makeLikeRequest(map[string]string{
		"post_id": "1",
		"value":   "1",
	}))

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestLike_NewReaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	h := &Handler{
		DB:   db,
		Auth: &mockAuth{ok: true, userID: 1},
	}

	// no existing row
	mock.ExpectQuery("SELECT value FROM reactions").
		WithArgs(1, 1).
		WillReturnError(sql.ErrNoRows)

	// insert new reaction
	mock.ExpectExec("INSERT INTO reactions").
		WithArgs(1, 1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	w := httptest.NewRecorder()
	h.Like(w, makeLikeRequest(map[string]string{
		"post_id": "1",
		"value":   "1",
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

func TestLike_ToggleOff(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	h := &Handler{
		DB:   db,
		Auth: &mockAuth{ok: true, userID: 1},
	}

	// existing row with same value — should delete (toggle off)
	mock.ExpectQuery("SELECT value FROM reactions").
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"value"}).AddRow(1))

	mock.ExpectExec("DELETE FROM reactions").
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	w := httptest.NewRecorder()
	h.Like(w, makeLikeRequest(map[string]string{
		"post_id": "1",
		"value":   "1",
	}))

	if w.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", w.Code)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled mock expectations: %v", err)
	}
}

func TestLike_SwitchReaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	h := &Handler{
		DB:   db,
		Auth: &mockAuth{ok: true, userID: 1},
	}

	// existing row with different value — should update (like -> dislike)
	mock.ExpectQuery("SELECT value FROM reactions").
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"value"}).AddRow(1))

	mock.ExpectExec("UPDATE reactions").
		WithArgs(-1, 1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	w := httptest.NewRecorder()
	h.Like(w, makeLikeRequest(map[string]string{
		"post_id": "1",
		"value":   "-1",
	}))

	if w.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", w.Code)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled mock expectations: %v", err)
	}
}

func TestLike_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	h := &Handler{
		DB:   db,
		Auth: &mockAuth{ok: true, userID: 1},
	}

	mock.ExpectQuery("SELECT value FROM reactions").
		WithArgs(1, 1).
		WillReturnError(sql.ErrConnDone)

	w := httptest.NewRecorder()
	h.Like(w, makeLikeRequest(map[string]string{
		"post_id": "1",
		"value":   "1",
	}))

	// handler logs the error but still redirects
	if w.Code != http.StatusSeeOther {
		t.Errorf("expected 303 (handler always redirects), got %d", w.Code)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled mock expectations: %v", err)
	}
}
