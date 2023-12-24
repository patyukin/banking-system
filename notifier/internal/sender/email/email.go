package email

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/patyukin/banking-system/notifier/internal/config/env"
	"log"
	"net/smtp"
)

var _ Sender = (*sender)(nil)

type Sender interface {
	Send(ctx context.Context, key string, value string) error
}

type sender struct {
	providers []env.EmailProvider
}

type KafkaValue struct {
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func (s *sender) Send(_ context.Context, _ string, value string) error {
	smtpHost := "fakesmtp"
	smtpPort := 1025
	from := "sender@example.com"

	var val KafkaValue
	err := json.Unmarshal([]byte(value), &val)
	if err != nil {
		return err
	}

	message := []byte("To: " + val.Email + "\r\n" +
		"Subject: " + val.Subject + "\r\n" +
		"\r\n" +
		val.Body + "\r\n")

	auth := smtp.PlainAuth("", "", "", smtpHost)
	err = smtp.SendMail(smtpHost+":"+string(rune(smtpPort)), auth, from, []string{val.Email}, message)
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
