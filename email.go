package main

import (
	// "log"

	"bytes"
	"html/template"
	"log"

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

func sendWelcomeEmail(to *User) error {
	em := &EmailMessage{}
	em.Subject = "Welcome to Anonutopia"
	em.FromName = "Anonutopia"
	em.FromEmail = "no-reply@anonutopia.com"
	em.BodyText = "Welcome to Anonutopia!"
	em.BodyHTML = "Welcome to Anonutopia!"
	em.ToEmail = to.Email
	em.ToName = to.Nickname

	t := template.New("welcome.html")
	var err error
	t, err = t.ParseFiles("emails/welcome.html")
	if err != nil {
		log.Println("fdsafad")
		return err
	}

	data := struct {
		VerificationLink string
	}{
		VerificationLink: "https://www.google.com",
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return err
	}

	em.BodyHTML = tpl.String()

	err = sendEmail(em)

	return err
}
