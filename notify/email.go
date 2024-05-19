package notify

import (
	"database/sql"
	"gopkg.in/gomail.v2"
	"log"
)

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
}

func SendEmailNotifications(db *sql.DB, message string, config EmailConfig) {
	rows, err := db.Query("SELECT email FROM subscription")
	if err != nil {
		log.Printf("Failed to retrieve subscriptions: %v\n", err)
		return
	}
	defer rows.Close()

	var emails []string
	for rows.Next() {
		var email string
		err := rows.Scan(&email)
		if err != nil {
			log.Printf("Failed to scan email: %v\n", err)
			continue
		}
		emails = append(emails, email)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", config.FromEmail)
	m.SetHeader("Subject", "Currency Rate Update")
	m.SetBody("text/plain", message)

	d := gomail.NewDialer(config.SMTPHost, config.SMTPPort, config.SMTPUser, config.SMTPPassword)

	for _, email := range emails {
		m.SetHeader("To", email)
		if err := d.DialAndSend(m); err != nil {
			log.Printf("Failed to send email to %s: %v\n", email, err)
		}
	}
}
