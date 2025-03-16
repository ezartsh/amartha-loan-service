package config

const (
	envVarMailHost        = "MAIL_HOST"
	envVarMailPort        = "MAIL_PORT"
	envVarMailUsername    = "MAIL_USERNAME"
	envVarMailPassword    = "MAIL_PASSWORD"
	envVarMailEmailSender = "MAIL_EMAIL_SENDER"
)

type Mail struct {
	host        string
	port        int
	username    string
	password    string
	emailSender string
}

func newMailConfig() Mail {
	return Mail{
		host:        getEnv(envVarMailHost, ""),
		port:        getEnvAsInt(envVarMailPort, 0),
		username:    getEnv(envVarMailUsername, ""),
		password:    getEnv(envVarMailPassword, ""),
		emailSender: getEnv(envVarMailEmailSender, ""),
	}
}

func (c Mail) Host() string {
	return c.host
}

func (c Mail) Port() int {
	return c.port
}

func (c Mail) Username() string {
	return c.username
}

func (c Mail) Password() string {
	return c.password
}

func (c Mail) EmailSender() string {
	return c.emailSender
}

func (c Mail) ValidateRequired(errs *[]EnvValidationError) {
	if c.host == "" {
		*errs = append(*errs, EnvValidationError{
			VariableName: envVarMailHost,
			Message:      "mail host is required",
			Description:  "e.g localhost",
		})
	}
	if c.port == 0 {
		*errs = append(*errs, EnvValidationError{
			VariableName: envVarMailPort,
			Message:      "mail port is required",
			Description:  "e.g 3000",
		})
	}
	if c.username == "" {
		*errs = append(*errs, EnvValidationError{
			VariableName: envVarMailUsername,
			Message:      "mail username is required",
			Description:  "e.g admin",
		})
	}
	if c.password == "" {
		*errs = append(*errs, EnvValidationError{
			VariableName: envVarMailPassword,
			Message:      "mail password is required",
			Description:  "e.g admin",
		})
	}
	if c.emailSender == "" {
		*errs = append(*errs, EnvValidationError{
			VariableName: envVarMailEmailSender,
			Message:      "mail email sender is required",
			Description:  "e.g admin@company.com",
		})
	}
}
