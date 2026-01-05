package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

const sessionTimeout = 5 * time.Minute

func InitDatabase(dbPath string) {

	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS session_tokens")
	if err != nil {
		log.Fatalf("Error dropping session_tokens table: %v", err)
	}

	createTables := []string{
		`CREATE TABLE IF NOT EXISTS accounts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			first_name TEXT NOT NULL DEFAULT 'Missing',
			last_name TEXT NOT NULL DEFAULT 'Missing',
			email TEXT NOT NULL UNIQUE,
			role TEXT NOT NULL DEFAULT 'client' CHECK(role IN ('client', 'admin'))
		)`,
		`CREATE TABLE IF NOT EXISTS session_tokens (
			token TEXT PRIMARY KEY,
			username TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			expires_at TIMESTAMP NOT NULL,  
			FOREIGN KEY(username) REFERENCES accounts(username)
		)`,
		`CREATE TABLE IF NOT EXISTS reservations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			table_number INTEGER NOT NULL,
			reservation_date TEXT NOT NULL,
			reservation_time TEXT NOT NULL,
			guests INTEGER NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending',
			email TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS tables (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			seats INTEGER NOT NULL,
			status TEXT NOT NULL DEFAULT 'available'
		)`,
	}

	for _, query := range createTables {
		_, err = db.Exec(query)
		if err != nil {
			log.Fatalf("Error creating table: %v", err)
		}
	}
}

// Close Database
func CloseDatabase() {
	if err := db.Close(); err != nil {
		log.Fatalf("Error closing database: %v", err)
	}
}

// Insert table
func InsertTable(seats int) error {

	_, err := db.Exec("INSERT INTO tables (seats) VALUES (?)", seats)
	if err != nil {
		log.Print("errore")
	}
	return nil
}
