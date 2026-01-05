package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"progetto/restaurant/server/database"

	"github.com/google/uuid"
)

var templates *template.Template

func init() {
	templates = template.Must(template.ParseGlob("server/templates/*.html"))
}

// Generate a session token
func generateSessionToken() (string, error) {
	token, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return token.String(), nil
}

func ValidateSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	_, err = database.ValidateSessionToken(cookie.Value)
	if err != nil {
		if err.Error() == "session expired" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}

// Middleware to check if user is admin
func RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := getUsernameFromSession(r)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		role, err := database.GetUserRole(username)
		if err != nil || role != "admin" {
			http.Error(w, "Forbidden: Admin access required", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}

// Middleware to check if user is client
func RequireClient(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := getUsernameFromSession(r)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		role, err := database.GetUserRole(username)
		if err != nil || role != "client" {
			http.Error(w, "Forbidden: Client access required", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}

func getUsernameFromSession(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return "", err
	}

	username, err := database.ValidateSessionToken(cookie.Value)
	if err != nil {
		if err.Error() == "session expired" {
			return "", fmt.Errorf("session expired")
		}
		return "", err
	}
	return username, nil
}
