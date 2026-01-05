package handler

import (
	"log"
	"net/http"
	"progetto/restaurant/server/database"
)

func InformationHandler(w http.ResponseWriter, r *http.Request) {
	ValidateSession(w, r)
	userInformation := Data{}

	username, err := getUsernameFromSession(r)
	if err != nil {
		http.Error(w, "Error retrieving username from session", http.StatusInternalServerError)
		return
	}

	// Manage the GET request
	if r.Method == http.MethodGet {
		firstName, lastName, email, err := database.GetUserInformation(username)
		if err != nil {
			userInformation.Error = "Missing information"
			err = templates.ExecuteTemplate(w, "account.html", userInformation)
			if err != nil {
				http.Error(w, "Error rendering template", http.StatusInternalServerError)
				return
			}
			return
		}

		if firstName == "Missing" {
			firstName = ""
		}
		if lastName == "Missing" {
			lastName = ""
		}

		userInformation = Data{
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
		}

		err = templates.ExecuteTemplate(w, "account.html", userInformation)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
		return

	}
	// Manage the POST request
	if r.Method == http.MethodPost {

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		userInformation.FirstName = r.FormValue("first_name")
		userInformation.LastName = r.FormValue("last_name")
		userInformation.Email = r.FormValue("email")

		// Insert the informations
		err = database.UpdateInformation(username, userInformation.FirstName, userInformation.LastName, userInformation.Email)
		if err != nil {
			log.Printf("Error inserting account: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}

func DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	ValidateSession(w, r)

	// Delete account
	if r.Method == http.MethodPost {
		log.Printf("metodo delete in corso")

		username, err := getUsernameFromSession(r)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Delete the user from the database
		if err := database.DeleteUser(username); err != nil {
			http.Error(w, "Error deleting account", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
