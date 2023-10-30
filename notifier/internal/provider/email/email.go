package email

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Provider struct {
	Name         string `yaml:"name"`
	SmtpPort     string `yaml:"port"`
	SmtpHost     string `yaml:"host"`
	SmtpLogin    string `yaml:"login"`
	SmtpPassword string `yaml:"password"`
	SmtpSecure   string `yaml:"secure"`
	MailFrom     string `yaml:"mailFrom"`
	NameFrom     string `yaml:"nameFrom"`
}

func New() (*[]Provider, error) {
	// get config sms providers list from email.yaml
	emailProvidersData, err := os.ReadFile("internal/config/email.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read email.yaml: %w", err)
	}

	// Parse the YAML data into a struct
	var emailProviders []Provider
	err = yaml.Unmarshal(emailProvidersData, &emailProviders)
	if err != nil {
		return nil, fmt.Errorf("failed to parse email.yaml: %w", err)
	}

	return &emailProviders, nil
}
