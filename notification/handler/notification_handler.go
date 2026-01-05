package handler

import (
	"encoding/json"
	"net/http"
	"progetto/notification/util"
)

type Notification struct {
	Recipient string `json:"recipient"`
	Subject   string `json:"subject"`
	Body      string `json:"message"`
}

func NotificationHandler(w http.ResponseWriter, r *http.Request) {
	var notif Notification

	if err := json.NewDecoder(r.Body).Decode(&notif); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if notif.Recipient == "" || notif.Subject == "" || notif.Body == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	util.SendEmail(notif.Recipient, notif.Subject, notif.Body)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email sent successfully"))
}
