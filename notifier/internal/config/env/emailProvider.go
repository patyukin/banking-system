package env

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type EmailProviderInterface interface {
	GetAll() []EmailProvider
	GetOneByName(name string) *EmailProvider
	GetDefault() *EmailProvider
}

type EmailProviderImpl struct {
	Providers []EmailProvider `json:"providers"`
}

type EmailProvider struct {
	Name         string
	SmtpPort     int8
	SmtpHost     string
	SmtpLogin    string
	SmtpPassword string
	SmtpSecure   bool
	MailFrom     string
	NameFrom     string
}

func (e *EmailProviderImpl) GetAll() ([]EmailProvider, error) {
	return e.Providers, nil
}

func (e *EmailProviderImpl) GetOneByName(name string) (*EmailProvider, error) {
	for _, provider := range e.Providers {
		if provider.Name == name {
			return &provider, nil
		}
	}

	return nil, fmt.Errorf("provider %s not found", name)
}

func (e *EmailProviderImpl) GetDefault() (*EmailProvider, error) {
	defaultProvider := os.Getenv("DEFAULT_EMAIL_PROVIDER")
	if len(defaultProvider) == 0 {
		return nil, fmt.Errorf("default provider not found")
	}

	for _, provider := range e.Providers {
		if provider.Name == defaultProvider {
			return &provider, nil
		}
	}

	return nil, fmt.Errorf("default provider not found")
}

func NewEmailProvider() ([]EmailProvider, error) {
	emailProviders := os.Getenv("EMAIL_PROVIDERS")
	if len(emailProviders) == 0 {
		return nil, fmt.Errorf("email Providers not found")
	}

	var emailProvidersData []EmailProvider
	err := yaml.Unmarshal([]byte(emailProviders), &emailProvidersData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse email Providers: %w", err)
	}

	return emailProvidersData, nil
}
