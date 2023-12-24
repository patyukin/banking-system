package env

import (
	"encoding/json"
	"fmt"
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
	Name         string `json:"name"`
	SmtpPort     int32  `json:"port"`
	SmtpHost     string `json:"host"`
	SmtpLogin    string `json:"login"`
	SmtpPassword string `json:"password"`
	SmtpSecure   bool   `json:"secure"`
	MailFrom     string `json:"mail_from"`
	NameFrom     string `json:"name_from"`
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
	err := json.Unmarshal([]byte(emailProviders), &emailProvidersData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse email Providers: %w", err)
	}

	return emailProvidersData, nil
}
