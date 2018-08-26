package email

import (
	"log"
	"net/smtp"
	"os"
)

var (
	auth smtp.Auth
	from string
	// Identity is the name of the sender
	Identity = "PRICE-TRACKER"
	username string
	password string
	host     string
	smtpPort = "587"
)

func init() {
	v, ok := os.LookupEnv("EMAIL_IDENTITY")
	if ok {
		Identity = v
	}

	v, ok = os.LookupEnv("EMAIL_USERNAME")
	if !ok {
		log.Println("WARN: EMAIL_USERNAME is not given in ENV")
	}
	username = v

	v, ok = os.LookupEnv("EMAIL_PASSWORD")
	if !ok {
		log.Println("WARN: EMAIL_PASSWORD is not given in ENV")
	}
	password = v

	v, ok = os.LookupEnv("EMAIL_HOST")
	if !ok {
		log.Println("WARN: EMAIL_HOST is not given in ENV")
	}
	host = v

	v, ok = os.LookupEnv("EMAIL_FROM")
	if ok {
		from = v
	}

	v, ok = os.LookupEnv("EMAIL_SMTP_PORT")
	if ok {
		smtpPort = v
	}

	auth = smtp.PlainAuth(
		Identity,
		username,
		password,
		host,
	)
}

// Send sends an email to destination
func Send(body string, subject string, destEmail string) (err error) {
	err = smtp.SendMail(host+":"+smtpPort, auth, from, []string{destEmail}, []byte(body))
	if err != nil {
		log.Printf("ERROR: SMTP error %s\n", err)
	}
	return
}
