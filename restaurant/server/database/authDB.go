package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Register a new user with role
func RegisterUser(username, hashedPassword, email, role string) error {
	_, err := db.Exec("INSERT INTO accounts (username, password, email, role) VALUES (?, ?, ?, ?)",
		username, hashedPassword, email, role)
	return err
}

// Remove a user
func DeleteUser(username string) error {
	_, err := db.Exec("DELETE FROM accounts WHERE username = ?", username)
	if err != nil {
		log.Printf("Error deleting user from database: %v", err)
		return err
	}
	return nil
}

// Get user's password by username (backward compatibility)
func GetUserPassword(username string) (string, error) {
	var password string
	err := db.QueryRow("SELECT password FROM accounts WHERE username = ?", username).Scan(&password)
	return password, err
}

// Get user's role
func GetUserRole(username string) (string, error) {
	var role string
	err := db.QueryRow("SELECT role FROM accounts WHERE username = ?", username).Scan(&role)
	return role, err
}

// Save session token
func SaveSessionToken(username, token string) error {
	expiresAt := time.Now().Add(sessionTimeout)
	_, err := db.Exec("INSERT INTO session_tokens (token, username, created_at, expires_at) VALUES (?, ?, ?, ?)",
		token, username, time.Now(), expiresAt)
	return err
}

// Validate session token
func ValidateSessionToken(token string) (string, error) {
	var username string
	var expiresAt time.Time
	err := db.QueryRow("SELECT username, expires_at FROM session_tokens WHERE token = ?", token).Scan(&username, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", err
		}
		return "", err
	}

	if time.Now().After(expiresAt) {
		return "", fmt.Errorf("session expired")
	}

	return username, nil
}

// Delete session token
func DeleteSessionToken(token string) error {
	_, err := db.Exec("DELETE FROM session_tokens WHERE token = ?", token)
	return err
}
