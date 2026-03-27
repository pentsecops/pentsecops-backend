package utils

import (
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

// SendEmail sends an email using SMTP
func SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_FROM"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		mustAtoi(os.Getenv("SMTP_PORT")),
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASS"),
	)
	return d.DialAndSend(m)
}

func mustAtoi(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}
