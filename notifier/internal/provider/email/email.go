package email

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
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
	// get config sms providers list from emailProviders.yaml
	emailProvidersData, err := os.ReadFile("internal/config/emailProviders.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read emailProviders.yaml: %w", err)
	}

	// Parse the YAML data into a struct
	var emailProviders []Provider
	err = yaml.Unmarshal(emailProvidersData, &emailProviders)
	if err != nil {
		return nil, fmt.Errorf("failed to parse emailProviders.yaml: %w", err)
	}

	return &emailProviders, nil
}
