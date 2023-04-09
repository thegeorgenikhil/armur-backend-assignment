package mail

import (
	"net/smtp"
)

// SendMailWithHTML sends an email with the given subject and html body to the given email address
func SendMailWithHTML(subject string, html string, from string, to []string, password string) error {
	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")

	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"

	msg := "Subject: " + subject + "\n" + headers + "\n\n" + html

	err := smtp.SendMail("smtp.gmail.com:587", auth, from, to, []byte(msg))

	return err
}
