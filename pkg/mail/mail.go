package mail

type MailInter interface {
	Send(to string, subject string, content string) error
}

type Mail struct {
	driver MailInter
}

func NewMail(driver MailInter) *Mail {
	return &Mail{
		driver: driver,
	}
}

func (m *Mail) Send(to string, subject string, content string) error {
	return m.driver.Send(to, subject, content)
}
