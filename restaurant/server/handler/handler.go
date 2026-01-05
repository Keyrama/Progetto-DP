package handler

import (
	"log"
	"net/http"
	"progetto/restaurant/server/database"

	"golang.org/x/crypto/bcrypt"
)

var userInformation Data

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Manage the GET request
	if r.Method == http.MethodGet {
		err := templates.ExecuteTemplate(w, "register.html", nil)
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

		userInformation.UserName = r.FormValue("username")
		userInformation.Password = r.FormValue("password")
		userInformation.Email = r.FormValue("email")

		// Password hashing
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInformation.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Register the user as 'client' by default
		err = database.RegisterUser(userInformation.UserName, string(hashedPassword), userInformation.Email, "client")
		if err != nil {
			log.Printf("Error registering user: %v", err)
			userInformation.Error = "Username already exists"
			templates.ExecuteTemplate(w, "register.html", userInformation)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// Manage the login request with role-based redirect
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := templates.ExecuteTemplate(w, "login.html", nil)
		if err != nil {
			http.Error(w, "Error rendering login page", http.StatusInternalServerError)
			return
		}
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Get password and role together
		hashedPassword, role, err := database.GetUserCredentials(username)
		if err != nil {
			userInformation.Error = "Invalid username"
			templates.ExecuteTemplate(w, "login.html", userInformation)
			return
		}

		// Password verification
		if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
			userInformation.Error = "Invalid password"
			templates.ExecuteTemplate(w, "login.html", userInformation)
			return
		}

		// Generate a session token
		sessionToken, err := generateSessionToken()
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Save the session token in the database
		if err := database.SaveSessionToken(username, sessionToken); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Set the session token in a cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    sessionToken,
			HttpOnly: true,
		})

		// Redirect based on role
		if role == "admin" {
			http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/home", http.StatusSeeOther)
		}
	}
}

// Manage the logout request
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ValidateSession(w, r)

	if r.Method == http.MethodGet {
		// Get the session token from the cookie
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "Session token missing or invalid", http.StatusUnauthorized)
			return
		}

		if err := database.DeleteSessionToken(cookie.Value); err != nil {
			http.Error(w, "Failed to log out", http.StatusInternalServerError)
			return
		}

		// Set the session token cookie to an empty value
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    "",
			HttpOnly: true,
			MaxAge:   -1,
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	ValidateSession(w, r)

	if r.Method == http.MethodGet {
		err := templates.ExecuteTemplate(w, "home.html", nil)
		if err != nil {
			http.Error(w, "Error rendering home page", http.StatusInternalServerError)
			return
		}
	}
}
