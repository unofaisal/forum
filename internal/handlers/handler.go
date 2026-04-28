package handlers

import (
	"database/sql"
	"net/http"
)

type AuthService interface {
	GetUserIDFromSession(r *http.Request) (int, bool)
}

type Handler struct {
	DB   *sql.DB
	Auth AuthService
}
