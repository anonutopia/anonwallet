package main

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	ui18n "github.com/Unknwon/i18n"
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var domains = [...]string{
	"aol.com", "att.net", "comcast.net", "facebook.com", "gmail.com", "gmx.com", "googlemail.com",
	"google.com", "hotmail.com", "hotmail.co.uk", "mac.com", "me.com", "mail.com", "msn.com",
	"live.com", "sbcglobal.net", "verizon.net", "yahoo.com", "yahoo.co.uk",

	"email.com", "fastmail.fm", "games.com", "gmx.net", "hush.com", "hushmail.com", "icloud.com",
	"iname.com", "inbox.com", "lavabit.com", "love.com", "outlook.com", "pobox.com", "protonmail.com",
	"rocketmail.com", "safe-mail.net", "wow.com", "ygm.com",
	"ymail.com", "zoho.com", "yandex.com",

	"bellsouth.net", "charter.net", "cox.net", "earthlink.net", "juno.com",

	"btinternet.com", "virginmedia.com", "blueyonder.co.uk", "freeserve.co.uk", "live.co.uk",
	"ntlworld.com", "o2.co.uk", "orange.net", "sky.com", "talktalk.co.uk", "tiscali.co.uk",
	"virgin.net", "wanadoo.co.uk", "bt.com",

	"sina.com", "sina.cn", "qq.com", "naver.com", "hanmail.net", "daum.net", "nate.com", "yahoo.co.jp",
	"yahoo.co.kr", "yahoo.co.id", "yahoo.co.in", "yahoo.com.sg", "yahoo.com.ph", "163.com", "126.com",
	"aliyun.com", "foxmail.com",

	"hotmail.fr", "live.fr", "laposte.net", "yahoo.fr", "wanadoo.fr", "orange.fr", "gmx.fr", "sfr.fr",
	"neuf.fr", "free.fr",

	"gmx.de", "hotmail.de", "live.de", "online.de", "t-online.de", "web.de", "yahoo.de",

	"libero.it", "virgilio.it", "hotmail.it", "aol.it", "tiscali.it", "alice.it", "live.it", "yahoo.it",
	"email.it", "tin.it", "poste.it", "teletu.it",

	"mail.ru", "rambler.ru", "yandex.ru", "ya.ru", "list.ru",

	"hotmail.be", "live.be", "skynet.be", "voo.be", "tvcablenet.be", "telenet.be",

	"hotmail.com.ar", "live.com.ar", "yahoo.com.ar", "fibertel.com.ar", "speedy.com.ar", "arnet.com.ar",

	"yahoo.com.mx", "live.com.mx", "hotmail.es", "hotmail.com.mx", "prodigy.net.mx",

	"yahoo.com.br", "hotmail.com.br", "outlook.com.br", "uol.com.br", "bol.com.br", "terra.com.br",
	"ig.com.br", "itelefonica.com.br", "r7.com", "zipmail.com.br", "globo.com", "globomail.com", "oi.com.br"}

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

func sendWelcomeEmail(to *User, lang string) error {
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
		return err
	}

	uid, err := encrypt([]byte(conf.DbPass[:16]), to.Address)
	if err != nil {
		return err
	}

	verLink := fmt.Sprintf(ui18n.Tr(lang, "verificationLink"), uid)

	data := struct {
		VerificationLink string
	}{
		VerificationLink: verLink,
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return err
	}

	em.BodyHTML = tpl.String()

	err = sendEmail(em)

	return err
}

func sendInitWelcomeEmail(to *User, password string, seed string, lang string) error {
	em := &EmailMessage{}
	em.Subject = "Welcome to Anonutopia"
	em.FromName = "Anonutopia"
	em.FromEmail = "no-reply@anonutopia.com"
	em.BodyText = "Welcome to Anonutopia!"
	em.BodyHTML = "Welcome to Anonutopia!"
	em.ToEmail = to.Email
	em.ToName = to.Nickname

	t := template.New("welcomeinit.html")
	var err error
	t, err = t.ParseFiles("emails/welcomeinit.html")
	if err != nil {
		return err
	}

	uid, err := encrypt([]byte(conf.DbPass[:16]), to.Email)
	if err != nil {
		return err
	}

	verLink := fmt.Sprintf(ui18n.Tr(lang, "verificationLink"), uid)

	data := struct {
		VerificationLink string
		Password         string
		Seed             string
	}{
		VerificationLink: verLink,
		Password:         password,
		Seed:             seed,
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return err
	}

	em.BodyHTML = tpl.String()

	err = sendEmail(em)

	return err
}

func sendPasswordResetEmail(to *User, lang string) error {
	em := &EmailMessage{}
	em.Subject = "Anonutopia Password Reset"
	em.FromName = "Anonutopia"
	em.FromEmail = "no-reply@anonutopia.com"
	em.BodyText = "Anonutopia Password Reset"
	em.BodyHTML = "Anonutopia Password Reset"
	em.ToEmail = to.Email
	em.ToName = to.Nickname

	t := template.New("passwordreset.html")
	var err error
	t, err = t.ParseFiles("emails/passwordreset.html")
	if err != nil {
		return err
	}

	uid, err := encrypt([]byte(conf.DbPass[:16]), to.Email)
	if err != nil {
		return err
	}

	resetLink := fmt.Sprintf(ui18n.Tr(lang, "passwordResetLink"), uid)

	data := struct {
		ResetLink string
	}{
		ResetLink: resetLink,
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return err
	}

	em.BodyHTML = tpl.String()

	err = sendEmail(em)

	return err
}

func validateEmailDomain(email string) bool {
	for i := range domains {
		if strings.HasSuffix(email, domains[i]) {
			return true
		}
	}
	return false
}
