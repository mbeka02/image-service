package mailer

import (
	gomail "gopkg.in/mail.v2"
)

type Mailer struct {
	Dialer *gomail.Dialer
}

func NewMailer(host, password string) *Mailer {
	return &Mailer{
		Dialer: gomail.NewDialer(host, 587, "api", password),
	}
}

func (m *Mailer) SendEmail() error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", "hello@demomailtrap.com")
	msg.SetHeader("To", "antonymbeka@gmail.com")

	msg.SetHeader("Subject", "Test email")

	// Set email body
	msg.SetBody("text/plain", "This is the Test Body")

	return m.Dialer.DialAndSend(msg)
}
