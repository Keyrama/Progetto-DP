package router_mux

import (
	"html/template"
	"net/http"
	"progetto/restaurant/server/handler"

	"github.com/gorilla/mux"
)

// Global variable for templates
var Templates *template.Template

// Inizializza il router
func InitRouter() *mux.Router {
	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./server/static"))))

	// Public routes
	r.HandleFunc("/", handler.LoginHandler).Methods("GET", "POST")
	r.HandleFunc("/register", handler.RegisterHandler).Methods("GET", "POST")
	r.HandleFunc("/logout", handler.LogoutHandler).Methods("GET")

	// Client routes
	r.HandleFunc("/home", handler.RequireClient(handler.HomePageHandler)).Methods("GET", "POST")
	r.HandleFunc("/account", handler.RequireClient(handler.InformationHandler)).Methods("GET", "POST")
	r.HandleFunc("/delete", handler.RequireClient(handler.DeleteAccountHandler)).Methods("POST")

	// Booking routes (Option B - multi-step)
	r.HandleFunc("/booking", handler.RequireClient(handler.BookingPageHandler)).Methods("GET")
	r.HandleFunc("/booking/step1", handler.RequireClient(handler.BookingStep1Handler)).Methods("POST")
	r.HandleFunc("/booking/create", handler.RequireClient(handler.CreateBookingHandler)).Methods("POST")
	r.HandleFunc("/my-bookings", handler.RequireClient(handler.MyBookingsHandler)).Methods("GET")

	// Admin routes
	r.HandleFunc("/admin/dashboard", handler.RequireAdmin(handler.AdminDashboardHandler)).Methods("GET")
	r.HandleFunc("/admin/confirm", handler.RequireAdmin(handler.ConfirmReservationHandler)).Methods("POST")
	r.HandleFunc("/admin/reject", handler.RequireAdmin(handler.RejectReservationHandler)).Methods("POST")

	return r
}

func SetTemplates(t *template.Template) {
	Templates = t
}
