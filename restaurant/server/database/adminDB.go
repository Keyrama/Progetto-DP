package database

import (
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Confirm a reservation
func ConfirmReservation(reservationID int) error {
	_, err := db.Exec("UPDATE reservations SET status = 'confirmed' WHERE id = ?", reservationID)
	return err
}

// Reject a reservation
func RejectReservation(reservationID int) error {
	_, err := db.Exec("UPDATE reservations SET status = 'rejected' WHERE id = ?", reservationID)
	return err
}

// Get count of today's reservations
func GetTodayReservationsCount() (int, error) {
	today := time.Now().Format("2006-01-02")
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM reservations WHERE reservation_date = ?", today).Scan(&count)
	return count, err
}

// Get count of pending reservations
func GetPendingReservationsCount() (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM reservations WHERE status = 'pending'").Scan(&count)
	return count, err
}

// Get count of available tables
func GetAvailableTablesCount() (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM tables WHERE status = 'available'").Scan(&count)
	return count, err
}

// Admin dashboard stats
type AdminStats struct {
	TodayReservations   int
	PendingReservations int
	AvailableTables     int
}

// Get all admin statistics
func GetAdminStats() (AdminStats, error) {
	var stats AdminStats
	var err error

	stats.TodayReservations, err = GetTodayReservationsCount()
	if err != nil {
		return stats, err
	}

	stats.PendingReservations, err = GetPendingReservationsCount()
	if err != nil {
		return stats, err
	}

	stats.AvailableTables, err = GetAvailableTablesCount()
	if err != nil {
		return stats, err
	}

	return stats, nil
}
