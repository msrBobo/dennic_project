package email

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"

	"dennic_api_gateway/internal/pkg/config"
)

type EmailData struct {
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Message     string `json:"message"`
}

func SendEmail(to []string, subject string, cfg config.Config, htmlpath string, data EmailData) error {
	t, err := template.ParseFiles(htmlpath)
	if err != nil {
		log.Println(err)
		return err
	}

	var k bytes.Buffer
	err = t.Execute(&k, data)
	if err != nil {
		return err
	}
	if k.String() == "" {
		fmt.Println("Error buffer")
	}
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte(fmt.Sprintf("Subject: %s\n%s\n%s", subject, mime, k.String()))
	// Authentication.
	auth := smtp.PlainAuth("", cfg.Email.SMTPEmail, cfg.Email.SMTPEmailPass, cfg.Email.SMTPHost)

	// Sending email.
	err = smtp.SendMail(cfg.Email.SMTPHost+":"+cfg.Email.SMTPPort, auth, cfg.Email.SMTPEmail, to, msg)
	return err
}
