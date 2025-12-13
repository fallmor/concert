package activity

import (
	"fmt"
	"net/smtp"
)

type EmailActivities struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
}

func NewEmailActivities(smtpHost string, smtpPort string, smtpUsername string, smtpPassword string) *EmailActivities {
	return &EmailActivities{
		SMTPHost:     smtpHost,
		SMTPPort:     smtpPort,
		SMTPUsername: smtpUsername,
		SMTPPassword: smtpPassword,
	}

}

func (e *EmailActivities) SendResetPasswordEmail(email string) error {
	auth := smtp.PlainAuth("", e.SMTPUsername, e.SMTPPassword, e.SMTPHost)
	from := "test@example.com"
	msg := []byte("To: " + email + "\r\n" +
		"Subject: Test email\r\n" +
		"\r\n" +
		"This a test\r\n")
	addr := fmt.Sprintf("%s:%s", e.SMTPHost, e.SMTPPort)

	err := smtp.SendMail(addr, auth, from, []string{email}, msg)
	if err != nil {
		return err
	}
	return nil
}
