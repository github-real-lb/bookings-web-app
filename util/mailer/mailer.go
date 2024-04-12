package mailer

import (
	"html/template"

	"github.com/github-real-lb/bookings-web-app/util/loggers"
	mail "github.com/xhit/go-simple-mail/v2"
)

// MailData holds an email
type MailData struct {
	To      string
	From    string
	Subject string
	Content template.HTML
}

// MailerChannel is a channel to pass emails data
type MailerChannel chan MailData

var mc = make(MailerChannel)

func GetMailerChannel() MailerChannel {
	return mc
}

func Listen(errChan loggers.ErrorChannel) {
	go func() {
		for {
			errChan <- sendMail(<-mc)
		}
	}()
}

func sendMail(m MailData) loggers.ErrorData {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	//server.Encryption = mail.EncryptionNone
	//server.Username = ""
	//server.Password = ""

	client, err := server.Connect()
	if err != nil {
		return loggers.ErrorData{
			Prefix: "error connecting to SMTP server to send mail",
			Error:  err,
		}
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).
		AddTo(m.To).
		SetSubject(m.Subject)

	email.SetBody(mail.TextHTML, "Hello, <strong>world</strong>!")

	return loggers.ErrorData{
		Prefix: "error sending mail",
		Error:  email.Send(client),
	}
}
