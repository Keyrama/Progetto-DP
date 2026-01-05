package util

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func SendEmail(recipient string, subject string, body string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading the .env file: %v", err)
	}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpEmail := os.Getenv("SMTP_EMAIL")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	if smtpHost == "" || smtpPort == "" || smtpEmail == "" || smtpPassword == "" {
		log.Fatal("One or more SMTP variables are not defined in the .env file")
	}

	port, err := strconv.Atoi(smtpPort)
	if err != nil {
		log.Fatalf("Error converting the SMTP port: %v", err)
	}

	message := gomail.NewMessage()
	message.SetHeader("From", smtpEmail)
	message.SetHeader("To", recipient)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", body)

	dialer := gomail.NewDialer(smtpHost, port, smtpEmail, smtpPassword)

	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println("Error sending the email:", err)
		panic(err)
	} else {
		fmt.Println("Email sent successfully!")
	}
}
