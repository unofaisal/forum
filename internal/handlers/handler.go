package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"forum/internal/auth"
)

type Handler struct {
	DB   *sql.DB
	Auth *auth.AuthHandler
}

func (h *Handler) RenderTemplate(w http.ResponseWriter, tmplName string, data interface{}) {
	tmpl, err := template.ParseFiles("ui/templates/base.html", "ui/templates/"+tmplName+".html")
	if err != nil {
		log.Printf("Error parsing templates: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
