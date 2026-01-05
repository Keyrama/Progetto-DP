package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"progetto/restaurant/server/database"
	"strconv"
)

type AdminDashboardData struct {
	Stats        database.AdminStats
	Reservations []database.Reservation
	Error        string
	Success      string
}

type EmailNotification struct {
	Recipient string `json:"recipient"`
	Subject   string `json:"subject"`
	Body      string `json:"message"`
}

// Helper function to send email notification
func sendEmailNotification(email, subject, body string) error {
	notification := EmailNotification{
		Recipient: email,
		Subject:   subject,
		Body:      body,
	}

	jsonData, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("error marshaling notification: %v", err)
	}

	resp, err := http.Post("http://localhost:8081/notification", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error sending notification: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("notification service returned status: %d", resp.StatusCode)
	}

	return nil
}

// Helper function to get reservation by ID
func getReservationByID(id int) (*database.Reservation, error) {
	reservations, err := database.GetAllReservations()
	if err != nil {
		return nil, err
	}

	for _, r := range reservations {
		if r.ID == id {
			return &r, nil
		}
	}

	return nil, fmt.Errorf("reservation not found")
}

// Admin Dashboard Handler - Display dashboard with stats and reservations
func AdminDashboardHandler(w http.ResponseWriter, r *http.Request) {
	ValidateSession(w, r)

	if r.Method == http.MethodGet {
		// Get admin statistics
		stats, err := database.GetAdminStats()
		if err != nil {
			log.Printf("Error getting admin stats: %v", err)
			http.Error(w, "Error loading dashboard", http.StatusInternalServerError)
			return
		}

		// Get all reservations
		reservations, err := database.GetAllReservations()
		if err != nil {
			log.Printf("Error getting reservations: %v", err)
			http.Error(w, "Error loading reservations", http.StatusInternalServerError)
			return
		}

		data := AdminDashboardData{
			Stats:        stats,
			Reservations: reservations,
		}

		err = templates.ExecuteTemplate(w, "adminDashboard.html", data)
		if err != nil {
			http.Error(w, "Error rendering dashboard", http.StatusInternalServerError)
			return
		}
	}
}

// Confirm Reservation Handler
func ConfirmReservationHandler(w http.ResponseWriter, r *http.Request) {
	ValidateSession(w, r)

	if r.Method == http.MethodPost {
		idStr := r.FormValue("reservation_id")
		if idStr == "" {
			http.Error(w, "Missing reservation ID", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid reservation ID", http.StatusBadRequest)
			return
		}

		// Get reservation details before confirming
		reservation, err := getReservationByID(id)
		if err != nil {
			log.Printf("Error getting reservation %d: %v", id, err)
			http.Error(w, "Error retrieving reservation", http.StatusInternalServerError)
			return
		}

		// Confirm the reservation
		err = database.ConfirmReservation(id)
		if err != nil {
			log.Printf("Error confirming reservation %d: %v", id, err)
			http.Error(w, "Error confirming reservation", http.StatusInternalServerError)
			return
		}

		// Send confirmation email
		subject := "Prenotazione Confermata - Crisbi's"
		body := fmt.Sprintf(`Gentile %s,

La tua prenotazione è stata confermata!

Dettagli della prenotazione:
- Data: %s
- Orario: %s
- Numero ospiti: %d
- Tavolo: %d

Ti aspettiamo da Crisbi's!

Cordiali saluti,
Il team di Crisbi's`,
			reservation.Name,
			reservation.ReservationDate,
			reservation.ReservationTime,
			reservation.Guests,
			reservation.TableNumber)

		err = sendEmailNotification(reservation.Email, subject, body)
		if err != nil {
			log.Printf("Warning: Failed to send confirmation email for reservation %d: %v", id, err)
			// Non blocchiamo l'operazione se l'email fallisce
		} else {
			log.Printf("Confirmation email sent successfully for reservation %d", id)
		}

		log.Printf("Reservation %d confirmed successfully", id)
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
	}
}

// Reject Reservation Handler
func RejectReservationHandler(w http.ResponseWriter, r *http.Request) {
	ValidateSession(w, r)

	if r.Method == http.MethodPost {
		idStr := r.FormValue("reservation_id")
		if idStr == "" {
			http.Error(w, "Missing reservation ID", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid reservation ID", http.StatusBadRequest)
			return
		}

		// Get reservation details before rejecting
		reservation, err := getReservationByID(id)
		if err != nil {
			log.Printf("Error getting reservation %d: %v", id, err)
			http.Error(w, "Error retrieving reservation", http.StatusInternalServerError)
			return
		}

		// Reject the reservation
		err = database.RejectReservation(id)
		if err != nil {
			log.Printf("Error rejecting reservation %d: %v", id, err)
			http.Error(w, "Error rejecting reservation", http.StatusInternalServerError)
			return
		}

		// Send rejection email
		subject := "Prenotazione Non Disponibile - Crisbi's"
		body := fmt.Sprintf(`Gentile %s,

Ci dispiace informarti che non è possibile confermare la tua prenotazione.

Dettagli della prenotazione richiesta:
- Data: %s
- Orario: %s
- Numero ospiti: %d

Motivo: Disponibilità esaurita per la data e l'orario richiesti.

Ti invitiamo a contattarci per trovare una soluzione alternativa o a effettuare una nuova prenotazione per un'altra data.

Ci scusiamo per l'inconveniente.

Cordiali saluti,
Il team di Crisbi's`,
			reservation.Name,
			reservation.ReservationDate,
			reservation.ReservationTime,
			reservation.Guests)

		err = sendEmailNotification(reservation.Email, subject, body)
		if err != nil {
			log.Printf("Warning: Failed to send rejection email for reservation %d: %v", id, err)
			// Non blocchiamo l'operazione se l'email fallisce
		} else {
			log.Printf("Rejection email sent successfully for reservation %d", id)
		}

		log.Printf("Reservation %d rejected successfully", id)
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
	}
}
