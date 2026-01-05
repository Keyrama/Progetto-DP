package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Create a new reservation
func CreateReservation(name, email string, tableNumber int, date, time string, guests int) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO reservations (name, email, table_number, reservation_date, reservation_time, guests, status)
		VALUES (?, ?, ?, ?, ?, ?, 'pending')`,
		name, email, tableNumber, date, time, guests)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// Cancel a reservation
func CancelReservation(reservationID int) error {
	_, err := db.Exec("UPDATE reservations SET status = 'canceled' WHERE id = ?", reservationID)
	return err
}

// Get all reservations (admin function)
func GetAllReservations() ([]Reservation, error) {
	rows, err := db.Query(`
		SELECT id, name, table_number, reservation_date, reservation_time, guests, status, email
		FROM reservations
		ORDER BY reservation_date DESC, reservation_time DESC
	`)
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

// Get available tables
func GetAvailableTables() ([]map[string]interface{}, error) {
	rows, err := db.Query("SELECT id, seats FROM tables WHERE status = 'available'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []map[string]interface{}
	for rows.Next() {
		var id, seats int
		err := rows.Scan(&id, &seats)

		if err != nil {
			return nil, err
		}

		tables = append(tables, map[string]interface{}{"id": id, "seats": seats})
	}
	return tables, nil
}

// Find available tables for a specific date and time
func FindAvailableTables(reservationDate, reservationTime string, guests int) (int, error) {
	startTime, err := time.Parse("15:04", reservationTime)
	if err != nil {
		return 0, fmt.Errorf("invalid time format: %v", err)
	}

	// Calculate end time (2 hours later)
	endTime := startTime.Add(2 * time.Hour)

	newStartMinutes := startTime.Hour()*60 + startTime.Minute()
	newEndMinutes := endTime.Hour()*60 + endTime.Minute()

	if newEndMinutes == 0 {
		newEndMinutes = 1440
	}

	query := `
		SELECT id, seats FROM tables
		WHERE seats >= ?
		AND status = 'available'
		AND id NOT IN (
			SELECT table_number FROM reservations
			WHERE reservation_date = ?
			AND status NOT IN ('canceled', 'rejected')
			AND (
				(CAST(substr(reservation_time, 1, 2) AS INTEGER) * 60 + CAST(substr(reservation_time, 4, 2) AS INTEGER)) < ?
				AND
				? < (CAST(substr(reservation_time, 1, 2) AS INTEGER) * 60 + CAST(substr(reservation_time, 4, 2) AS INTEGER) + 120)
			)
		)
		ORDER BY seats ASC
		LIMIT 1;
	`

	var tableID, seats int
	err = db.QueryRow(query, guests, reservationDate, newEndMinutes, newStartMinutes).Scan(&tableID, &seats)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("no available table found")
		}
		return 0, err
	}
	return tableID, nil
}

// Get available time slots for a specific date and number of guests
func GetAvailableTimeSlots(date string, guests int) ([]string, error) {
	lunchSlots := []string{"12:00", "12:30", "13:00", "13:30", "14:00", "14:30"}
	dinnerSlots := []string{"19:00", "19:30", "20:00", "20:30", "21:00", "21:30", "22:00"}

	allSlots := append(lunchSlots, dinnerSlots...)
	availableSlots := []string{}

	for _, timeSlot := range allSlots {
		tableID, err := FindAvailableTables(date, timeSlot, guests)
		if err == nil && tableID > 0 {
			availableSlots = append(availableSlots, timeSlot)
		}
	}

	return availableSlots, nil
}

// Reservation struct
type Reservation struct {
	ID              int
	Name            string
	TableNumber     int
	ReservationDate string
	ReservationTime string
	Guests          int
	Status          string
	Email           string
}
