package mailer

import (
	"github.com/mbeka02/image-service/config"
	gomail "gopkg.in/mail.v2"
)

func SendEmail() error {
	msg := gomail.NewMessage()
	conf, err := config.LoadConfig("../..")
	if err != nil {
		return err
	}
	msg.SetHeader("From", "hello@demomailtrap.com")
	msg.SetHeader("To", "antonymbeka@gmail.com")

	msg.SetHeader("Subject", "Test email")

	// Set email body
	msg.SetBody("text/plain", "This is the Test Body")

	dialer := gomail.NewDialer(conf.MAILER_HOST, 587, "api", conf.MAILER_PASSWORD)

	return dialer.DialAndSend(msg)
}
