package main

import (
	"log"
	"os"

	"PR3_CSA_Hunko/smtp"
)

func main() {
	config := smtp.Config{
		Host:     getenv("SMTP_HOST", "smtp.gmail.com"),
		Port:     getenv("SMTP_PORT", "587"),
		From:     os.Getenv("SMTP_FROM"),
		Password: os.Getenv("SMTP_PASSWORD"),
	}

	sender := smtp.NewSender(config)
	to := []string{getenv("SMTP_TEST_TO", "duty-engineer@internal.net")}

	if err := sender.SendEmail(to, "Test Email", "This is a test email sent from Go."); err != nil {
		log.Fatal(err)
	}
	log.Println("Test email sent to", to)
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
