package helpers

import (
	"log"

	"github.com/google/uuid"
	"gopkg.in/gomail.v2"
)

type EmailJob struct {
	EmailID uuid.UUID `json:"email_id"`
	To      string    `json:"to"`
	Subject string    `json:"subject"`
	Body    string    `json:"body"`
}

func SendEmail(job EmailJob, smtpHost string, smtpPort int, smtpUser, smtpPass string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", smtpUser)
	m.SetHeader("To", job.To)
	m.SetHeader("Subject", job.Subject)
	m.SetBody("text/html", job.Body)

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	if err := d.DialAndSend(m); err != nil {
		log.Println("Failed to send email:", err)
		return err
	} else {
		log.Println("Email sent to", job.To)
		return nil
	}
}
