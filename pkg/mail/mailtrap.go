package mail

import (
	"loan-service/config"

	gomail "gopkg.in/mail.v2"
)

type MailTrap struct {
	host        string
	port        int
	username    string
	password    string
	emailSender string
}

func NewMailTrap() MailTrap {
	return MailTrap{
		host:        config.Env.Mail.Host(),
		port:        config.Env.Mail.Port(),
		username:    config.Env.Mail.Username(),
		password:    config.Env.Mail.Password(),
		emailSender: config.Env.Mail.EmailSender(),
	}
}

func (m MailTrap) Send(to string, subject string, content string) error {
	message := gomail.NewMessage()

	// Set email headers
	message.SetHeader("From", m.emailSender)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)

	// Set email body
	message.SetBody("text/html", content)

	// Set up the SMTP dialer
	dialer := gomail.NewDialer(
		m.host,
		m.port,
		m.username,
		m.password,
	)

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		return err
	}

	return nil
}
