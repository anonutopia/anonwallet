package main

import (
	// "log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailMessage struct {
	FromName  string
	FromEmail string
	ToName    string
	ToEmail   string
	Subject   string
	BodyHTML  string
	BodyText  string
}

func sendEmail(em *EmailMessage) error {
	from := mail.NewEmail(em.FromName, em.FromEmail)
	to := mail.NewEmail(em.ToName, em.ToEmail)
	message := mail.NewSingleEmail(from, em.Subject, to, em.BodyText, em.BodyHTML)

	client := sendgrid.NewSendClient(conf.SendgridKey)
	_, err := client.Send(message)

	return err
}
