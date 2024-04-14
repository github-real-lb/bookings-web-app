package mailers

import (
	"html/template"
	"sync"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

// MailData holds an email
type MailData struct {
	To      string
	From    string
	Subject string
	Content template.HTML
}

type SmartMailer struct {
	*mail.SMTPServer               // SMTP Server
	MailChannel      chan MailData //channel to pass emails data

	done     chan struct{} // used to stop the ListenAndMail() function
	shutdown sync.Once     // ensures Shutdown() is only performed once
}

func NewSmartMailer() *SmartMailer {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	//server.Encryption = mail.EncryptionNone
	//server.Username = ""
	//server.Password = ""

	return &SmartMailer{
		SMTPServer: server,
	}

}

func (sm *SmartMailer) SendMail(data MailData) error {
	client, err := sm.Connect()
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(data.From).
		AddTo(data.To).
		SetSubject(data.Subject)

	email.SetBody(mail.TextHTML, string(data.Content))

	err = email.Send(client)

	if err != nil {
		return err
	}

	return nil
}

// ListenAndMail listens for MailData on MailerChannel and sends emails
// errChan is used to send errors
// buffer determine the buffer size of the channel. buffer = 100 is the minimum
// Make sure to use Shutdown() to stop listening and close channel
func (sm *SmartMailer) ListenAndMail(errChan chan any, buffer int) {
	if buffer < 100 {
		buffer = 100
	}

	// create mail channel with buffer size of 100
	sm.MailChannel = make(chan MailData, buffer)

	// create the done channel to stop the listening
	sm.done = make(chan struct{})

	var err error

	// start listening
	for {
		select {
		case v := <-sm.MailChannel:
			// sending mail
			err = sm.SendMail(v)
			if err != nil {
				errChan <- err
			}
		case <-sm.done:
			// wait for ensure channel complete sending mail
			time.Sleep(100 * time.Millisecond)

			// close channel
			close(sm.MailChannel)
			return
		}
	}
}

// Shutdown stops ListenAndMail() and close channels
func (sm *SmartMailer) Shutdown() {
	if sm.done == nil {
		return
	}

	// close done channel safely
	sm.shutdown.Do(func() {
		close(sm.done)
	})
}
