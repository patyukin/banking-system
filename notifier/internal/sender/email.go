package sender

import (
	"fmt"
	"github.com/patyukin/banking-system/notifier/internal/config/env"
	"log"
	"net/smtp"
)

var _ Sender = (*sender)(nil)

type Sender interface {
	Send(to string, text string) error
}

type sender struct {
	providers []env.EmailProvider
}

func (s *sender) Send(to string, text string) error {
	smtpHost := "fakesmtp"
	smtpPort := 1025
	from := "sender@example.com"

	subject := "Test Email"
	message := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		text + "\r\n")

	auth := smtp.PlainAuth("", "", "", smtpHost)
	err := smtp.SendMail(smtpHost+":"+string(rune(smtpPort)), auth, from, []string{to}, message)
	if err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}

	log.Println("Email has been sent successfully.")
	return nil
}

func New(providers []env.EmailProvider) (Sender, error) {
	return &sender{
		providers: providers,
	}, nil
}
