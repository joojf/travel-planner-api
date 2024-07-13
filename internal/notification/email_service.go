package notification

import (
	"fmt"
	"net/smtp"
)

type EmailService struct {
	smtpHost     string
	smtpPort     int
	smtpUsername string
	smtpPassword string
	fromEmail    string
}

func NewEmailService(host string, port int, username, password, fromEmail string) *EmailService {
	return &EmailService{
		smtpHost:     host,
		smtpPort:     port,
		smtpUsername: username,
		smtpPassword: password,
		fromEmail:    fromEmail,
	}
}

func (s *EmailService) SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.smtpUsername, s.smtpPassword, s.smtpHost)

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", s.fromEmail, to, subject, body)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", s.smtpHost, s.smtpPort),
		auth,
		s.fromEmail,
		[]string{to},
		[]byte(msg),
	)

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
