package handlers

import (
	"database/sql"

	"forum/internal/auth"
)

type Handler struct {
	DB   *sql.DB
	Auth *auth.AuthHandler
}
