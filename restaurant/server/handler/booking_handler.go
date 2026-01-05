package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"progetto/restaurant/server/database"
	"strconv"
	"time"
)

type BookingPageData struct {
	Error          string
	Success        string
	Date           string
	Guests         int
	AvailableTimes []string
	LunchTimes     []string
	DinnerTimes    []string
}

// Helper function to render booking page with data
func renderBookingPage(w http.ResponseWriter, data BookingPageData) {
	err := templates.ExecuteTemplate(w, "booking.html", data)
	if err != nil {
		log.Printf("Error rendering booking page: %v", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}

// Helper function to separate lunch and dinner times
func separateLunchDinner(times []string) ([]string, []string) {
	lunchTimes := []string{}
	dinnerTimes := []string{}

	for _, t := range times {
		hour := 0
		_, err := fmt.Sscanf(t, "%d:", &hour)
		if err == nil {
			if hour >= 12 && hour < 15 {
				lunchTimes = append(lunchTimes, t)
			} else if hour >= 19 {
				dinnerTimes = append(dinnerTimes, t)
			}
		}
	}

	return lunchTimes, dinnerTimes
}

// Booking page handler - Step 1: Show form for date and guests
func BookingPageHandler(w http.ResponseWriter, r *http.Request) {
	ValidateSession(w, r)

	if r.Method == http.MethodGet {
		renderBookingPage(w, BookingPageData{})
	}
}

// Booking Step 1 Handler - Process date and guests, show available times
func BookingStep1Handler(w http.ResponseWriter, r *http.Request) {
	ValidateSession(w, r)

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/booking", http.StatusSeeOther)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v", err)
		renderBookingPage(w, BookingPageData{Error: "Errore nel form. Riprova."})
		return
	}

	// Get form values
	date := r.FormValue("date")
	guestsStr := r.FormValue("guests")

	// Validate form values
	if date == "" || guestsStr == "" {
		renderBookingPage(w, BookingPageData{Error: "Tutti i campi sono obbligatori."})
		return
	}

	// Parse guests
	guests, err := strconv.Atoi(guestsStr)
	if err != nil || guests < 1 || guests > 6 {
		renderBookingPage(w, BookingPageData{Error: "Numero di ospiti non valido (1-6)."})
		return
	}

	// Validate date (not in the past)
	bookingDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		renderBookingPage(w, BookingPageData{Error: "Data non valida."})
		return
	}

	today := time.Now().Truncate(24 * time.Hour)
	if bookingDate.Before(today) {
		renderBookingPage(w, BookingPageData{
			Error:  "Non puoi prenotare per una data passata.",
			Date:   date,
			Guests: guests,
		})
		return
	}

	// Get available time slots
	availableTimes, err := database.GetAvailableTimeSlots(date, guests)
	if err != nil {
		log.Printf("Error getting available time slots: %v", err)
		renderBookingPage(w, BookingPageData{
			Error:  "Errore nel recupero degli orari disponibili.",
			Date:   date,
			Guests: guests,
		})
		return
	}

	// Filter out past times if booking is for today
	if bookingDate.Equal(today) {
		now := time.Now()
		var filteredTimes []string
		for _, t := range availableTimes {
			bookingTime, err := time.Parse("15:04", t)
			if err != nil {
				continue
			}
			bookingDateTime := time.Date(now.Year(), now.Month(), now.Day(),
				bookingTime.Hour(), bookingTime.Minute(), 0, 0, now.Location())

			if bookingDateTime.After(now) {
				filteredTimes = append(filteredTimes, t)
			}
		}
		availableTimes = filteredTimes
	}

	// Separate lunch and dinner times
	lunchTimes, dinnerTimes := separateLunchDinner(availableTimes)

	// Show step 2 with available times
	data := BookingPageData{
		Date:           date,
		Guests:         guests,
		AvailableTimes: availableTimes,
		LunchTimes:     lunchTimes,
		DinnerTimes:    dinnerTimes,
	}

	if len(availableTimes) == 0 {
		data.Error = "Nessun tavolo disponibile per questa data e numero di ospiti. Prova un'altra data."
	}

	renderBookingPage(w, data)
}

// Create booking handler - Final step: Create the reservation
func CreateBookingHandler(w http.ResponseWriter, r *http.Request) {
	ValidateSession(w, r)

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/booking", http.StatusSeeOther)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v", err)
		renderBookingPage(w, BookingPageData{Error: "Errore nel form. Riprova."})
		return
	}

	// Get form values
	date := r.FormValue("date")
	timeSlot := r.FormValue("time")
	guestsStr := r.FormValue("guests")

	// Validate form values
	if date == "" || timeSlot == "" || guestsStr == "" {
		renderBookingPage(w, BookingPageData{Error: "Tutti i campi sono obbligatori."})
		return
	}

	// Parse guests
	guests, err := strconv.Atoi(guestsStr)
	if err != nil || guests < 1 || guests > 6 {
		renderBookingPage(w, BookingPageData{Error: "Numero di ospiti non valido (1-6)."})
		return
	}

	// Validate date (not in the past)
	bookingDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		renderBookingPage(w, BookingPageData{Error: "Data non valida."})
		return
	}

	today := time.Now().Truncate(24 * time.Hour)
	now := time.Now()

	if bookingDate.Before(today) {
		renderBookingPage(w, BookingPageData{
			Error:  "Non puoi prenotare per una data passata.",
			Date:   date,
			Guests: guests,
		})
		return
	}

	// Validate time format
	bookingTime, err := time.Parse("15:04", timeSlot)
	if err != nil {
		renderBookingPage(w, BookingPageData{Error: "Orario non valido."})
		return
	}

	// Check if booking is for today and time has already passed
	if bookingDate.Equal(today) {
		bookingDateTime := time.Date(now.Year(), now.Month(), now.Day(),
			bookingTime.Hour(), bookingTime.Minute(), 0, 0, now.Location())

		if bookingDateTime.Before(now) {
			renderBookingPage(w, BookingPageData{
				Error: "Non puoi prenotare per un orario già passato.",
			})
			return
		}
	}

	// Check if time is within valid slots
	hourInt := bookingTime.Hour()
	if !((hourInt >= 12 && hourInt < 15) || (hourInt >= 19 && hourInt < 23)) {
		renderBookingPage(w, BookingPageData{
			Error: "Orario non valido. Scegli tra pranzo (12:00-14:30) o cena (19:00-22:00).",
		})
		return
	}

	// Get username from session
	username, err := getUsernameFromSession(r)
	if err != nil {
		renderBookingPage(w, BookingPageData{
			Error: "Errore di autenticazione. Effettua nuovamente il login.",
		})
		return
	}

	// Get user info
	firstName, lastName, email, err := database.GetUserInformation(username)
	if err != nil {
		log.Printf("Error retrieving user info: %v", err)
		renderBookingPage(w, BookingPageData{
			Error: "Errore nel recupero delle informazioni utente.",
		})
		return
	}

	// Find available table
	tableID, err := database.FindAvailableTables(date, timeSlot, guests)
	if err != nil {
		log.Printf("Error finding available table: %v", err)

		// Get available times again to show them
		availableTimes, _ := database.GetAvailableTimeSlots(date, guests)
		lunchTimes, dinnerTimes := separateLunchDinner(availableTimes)

		renderBookingPage(w, BookingPageData{
			Error:          "Questo orario non è più disponibile. Seleziona un altro orario.",
			Date:           date,
			Guests:         guests,
			AvailableTimes: availableTimes,
			LunchTimes:     lunchTimes,
			DinnerTimes:    dinnerTimes,
		})
		return
	}

	// Create reservation
	reservationID, err := database.CreateReservation(
		firstName+" "+lastName,
		email,
		tableID,
		date,
		timeSlot,
		guests,
	)
	if err != nil {
		log.Printf("Error creating reservation: %v", err)
		renderBookingPage(w, BookingPageData{
			Error: "Errore nella creazione della prenotazione. Riprova.",
		})
		return
	}

	log.Printf("Reservation created successfully: ID=%d, Table=%d", reservationID, tableID)
	renderBookingPage(w, BookingPageData{
		Success: "Prenotazione creata con successo! In attesa di conferma dall'amministratore.",
	})
}

// Get Time slots for a specific date (API endpoint - kept for compatibility)
func GetAvailableTimeSlotsHandler(w http.ResponseWriter, r *http.Request) {
	ValidateSession(w, r)

	if r.Method != http.MethodGet {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]interface{}{
			"success": false,
			"message": "Method not allowed",
		})
		return
	}

	date := r.URL.Query().Get("date")
	guestsStr := r.URL.Query().Get("guests")

	if date == "" || guestsStr == "" {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Missing date or guests parameter",
		})
		return
	}

	guests, err := strconv.Atoi(guestsStr)
	if err != nil || guests <= 0 {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid guests parameter",
		})
		return
	}

	// Get available time slots
	availableTimes, err := database.GetAvailableTimeSlots(date, guests)
	if err != nil {
		log.Printf("Error retrieving available time slots: %v", err)
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Error retrieving available time slots",
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"times":   availableTimes,
	})
}

// Handler for displaying user's bookings
func MyBookingsHandler(w http.ResponseWriter, r *http.Request) {
	ValidateSession(w, r)

	username, err := getUsernameFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Get user's email
	_, _, email, err := database.GetUserInformation(username)
	if err != nil {
		http.Error(w, "Error retrieving user information", http.StatusInternalServerError)
		return
	}

	// Get user's reservations
	reservations, err := database.GetUserReservations(email)
	if err != nil {
		log.Printf("Error getting user reservations: %v", err)
		http.Error(w, "Error retrieving reservations", http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(w, "myBookings.html", reservations)
	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}
}

// JSON response helper (for API endpoints)
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
