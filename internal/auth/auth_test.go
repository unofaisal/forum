package auth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func setupHandler(t *testing.T) (*AuthHandler, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	return &AuthHandler{DB: db}, mock
}

func makeRegisterRequest(fields map[string]string) *http.Request {
	form := url.Values{}
	for k, v := range fields {
		form.Set(k, v)
	}
	req := httptest.NewRequest("POST", "/register", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

func TestRegister_MissingFields(t *testing.T) {
	tests := []struct {
		name   string
		fields map[string]string
	}{
		{
			name: "missing username",
			fields: map[string]string{
				"email": "test@example.com", "password": "secret123", "confirmPassword": "secret123",
			},
		},
		{
			name: "missing email",
			fields: map[string]string{
				"username": "john", "password": "secret123", "confirmPassword": "secret123",
			},
		},
		{
			name: "missing password",
			fields: map[string]string{
				"username": "john", "email": "test@example.com", "confirmPassword": "secret123",
			},
		},
		{
			name: "missing confirmPassword",
			fields: map[string]string{
				"username": "john", "email": "test@example.com", "password": "secret123",
			},
		},
		{
			name:   "all fields missing",
			fields: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, _ := setupHandler(t)
			req := makeRegisterRequest(tt.fields)
			w := httptest.NewRecorder()

			h.Register(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("expected 400, got %d", w.Code)
			}
			if !strings.Contains(w.Body.String(), "all fields are required") {
				t.Errorf("expected error message, got: %s", w.Body.String())
			}
		})
	}
}

func TestRegister_PasswordMismatch(t *testing.T) {
	h, _ := setupHandler(t)

	req := makeRegisterRequest(map[string]string{
		"username":        "john",
		"email":           "john@example.com",
		"password":        "secret123",
		"confirmPassword": "different456",
	})
	w := httptest.NewRecorder()

	h.Register(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "passwords do not match") {
		t.Errorf("expected passwords do not match error, got: %s", w.Body.String())
	}
}

func TestRegister_DBError(t *testing.T) {
	h, mock := setupHandler(t)

	mock.ExpectExec("INSERT INTO users").
		WithArgs("john", sqlmock.AnyArg(), "john@example.com").
		WillReturnError(fmt.Errorf("UNIQUE constraint failed: users.username"))

	req := makeRegisterRequest(map[string]string{
		"username":        "john",
		"email":           "john@example.com",
		"password":        "secret123",
		"confirmPassword": "secret123",
	})
	w := httptest.NewRecorder()

	h.Register(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "username or email already exists") {
		t.Errorf("expected db error message, got: %s", w.Body.String())
	}
}

func TestRegister_Success(t *testing.T) {
	h, mock := setupHandler(t)

	mock.ExpectExec("INSERT INTO users").
		WithArgs("john", sqlmock.AnyArg(), "john@example.com").
		WillReturnResult(sqlmock.NewResult(1, 1))

	req := makeRegisterRequest(map[string]string{
		"username":        "john",
		"email":           "john@example.com",
		"password":        "secret123",
		"confirmPassword": "secret123",
	})
	w := httptest.NewRecorder()

	h.Register(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("expected 303 redirect, got %d", w.Code)
	}
	if w.Header().Get("Location") != "/" {
		t.Errorf("expected redirect to /, got: %s", w.Header().Get("Location"))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled mock expectations: %v", err)
	}
}
