package sms

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Provider struct {
	Name      string `yaml:"name"`
	Url       string `yaml:"url"`
	Login     string `yaml:"login"`
	Password  string `yaml:"password"`
	IsDefault bool   `yaml:"isDefault"`
}

func New() (*[]Provider, error) {
	// get config sms providers list from emailProviders.yaml
	smsProvidersData, err := os.ReadFile("internal/config/smsProviders.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read smsProviders.yaml: %w", err)
	}

	// Parse the YAML data into a struct
	var smsProviders []Provider
	err = yaml.Unmarshal(smsProvidersData, &smsProviders)
	if err != nil {
		return nil, fmt.Errorf("failed to parse smsProviders.yaml: %w", err)
	}

	return &smsProviders, nil
}

func (p *Provider) Get() (Provider, error) {

	return *p, nil
}
