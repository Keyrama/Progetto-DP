package database

import (
	_ "github.com/mattn/go-sqlite3"
)

// Update information
func UpdateInformation(userName, firstName, lastName, email string) error {
	_, err := db.Exec("UPDATE accounts SET first_name = ?, last_name = ?, email = ? WHERE username = ?",
		firstName, lastName, email, userName)
	return err
}

// Get user's information
func GetUserInformation(username string) (string, string, string, error) {
	var firstName, lastName, email string
	err := db.QueryRow("SELECT first_name, last_name, email FROM accounts WHERE username = ?", username).Scan(&firstName, &lastName, &email)
	if err != nil {
		return "", "", "", err
	}
	return firstName, lastName, email, nil
}

// Get user's password and role by username
func GetUserCredentials(username string) (string, string, error) {
	var password, role string
	err := db.QueryRow("SELECT password, role FROM accounts WHERE username = ?", username).Scan(&password, &role)
	return password, role, err
}

// Get user's reservations
func GetUserReservations(email string) ([]Reservation, error) {
	rows, err := db.Query(`
		SELECT id, name, table_number, reservation_date, reservation_time, guests, status, email
		FROM reservations
		WHERE email = ?
		ORDER BY reservation_date DESC, reservation_time DESC
	`, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []Reservation
	for rows.Next() {
		var r Reservation
		err := rows.Scan(&r.ID, &r.Name, &r.TableNumber, &r.ReservationDate, &r.ReservationTime, &r.Guests, &r.Status, &r.Email)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, r)
	}

	return reservations, nil
}
