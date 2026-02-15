package activity

import (
	"bytes"
	"concert/internal/utils"
	"fmt"
	"html/template"

	"github.com/resend/resend-go/v2"
)

type EmailActivities struct {
	Client resend.Client
	Token  string
}

func NewEmailActivities(token string) *EmailActivities {
	return &EmailActivities{
		Client: *resend.NewClient(token),
	}

}

func (e *EmailActivities) SendResetPasswordEmail(email, password string) error {
	tmpl, err := template.ParseFiles("mail.html")
	if err != nil {
		return err
	}

	appURL := utils.GetEnvOrDefault("APP_URL", "http://localhost:8080")
	inf := map[string]string{
		"Password": password,
		"URL":      appURL + "/login",
	}
	var body bytes.Buffer
	if err := tmpl.Execute(&body, inf); err != nil {
		return err
	}
	params := &resend.SendEmailRequest{
		From:    "Concert Booking system <onboarding@resend.dev>",
		To:      []string{email},
		Subject: "Password reset",
		Html:    body.String(),
	}

	sent, err := e.Client.Emails.Send(params)
	if err != nil {
		panic(err)
	}
	fmt.Println(sent.Id)
	return nil
}
