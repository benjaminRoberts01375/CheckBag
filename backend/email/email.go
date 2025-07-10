package email

import (
	Coms "github.com/benjaminRoberts01375/Go-Communicate"
	"github.com/resend/resend-go/v2"
)

func sendEmail(to string, subject string, message string) {
	if !Config.AllowSendingEmails {
		Coms.Println("Not sending email to " + to)
		return
	}

	client := resend.NewClient(Config.EmailAPIKey)
	emailParams := &resend.SendEmailRequest{
		From:    Config.EmailFrom,
		To:      []string{to},
		Subject: subject,
		Text:    message,
	}
	client.Emails.Send(emailParams)
}
