package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"progetto/restaurant/server/database"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	database.InitDatabase("../server/restaurant.db")
	defer database.CloseDatabase()

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== Create Admin User ===")
	fmt.Println()

	fmt.Print("Enter admin username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Enter admin password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	email := "crisbi.restaurant@gmail.com"

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Error hashing password: %v", err)
	}

	// Create admin user
	err = database.RegisterUser(username, string(hashedPassword), email, "admin")
	if err != nil {
		log.Fatalf("Error creating admin: %v", err)
	}

	fmt.Println()
	fmt.Println("Admin user created successfully!")
	fmt.Printf("Username: %s\n", username)
	fmt.Printf("Email: %s\n", email)
	fmt.Println("Role: admin")
	fmt.Println()
	fmt.Println("You can now login at http://localhost:8080")
}
