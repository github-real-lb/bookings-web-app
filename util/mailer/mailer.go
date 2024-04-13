package mailer

import (
	"html/template"

	"github.com/github-real-lb/bookings-web-app/util"
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

var mailerChan = make(MailerChannel)

func GetMailerChannel() MailerChannel {
	return mailerChan
}

func Listen(errChan chan error) {
	go func() {
		for {
			errChan <- sendMail(<-mailerChan)
		}
	}()
}

func sendMail(m MailData) error {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	//server.Encryption = mail.EncryptionNone
	//server.Username = ""
	//server.Password = ""

	client, err := server.Connect()
	if err != nil {
		return util.NewText().
			AddLineIndent("error connecting to SMTP server to send mail", "\t").
			AddLineIndent(err, "\t")
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).
		AddTo(m.To).
		SetSubject(m.Subject)

	email.SetBody(mail.TextHTML, "Hello, <strong>world</strong>!")

	err = email.Send(client)

	return util.NewText().
		AddLineIndent("error sending mail", "\t").
		AddLineIndent(err, "\t")
}
