package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ClientID     string
	ClientSecret string
	Password     string
	Subreddit    string
	UserAgent    string
	Username     string
}

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("could not load env file: %w", err)
	}

	config := Config{}
	err = config.GetStringFromENV(&config.ClientID, "CLIENT_ID")
	if err != nil {
		return nil, err
	}

	err = config.GetStringFromENV(&config.ClientSecret, "CLIENT_SECRET")
	if err != nil {
		return nil, err
	}

	err = config.GetStringFromENV(&config.Password, "PASSWORD")
	if err != nil {
		return nil, err
	}

	err = config.GetStringFromENV(&config.Subreddit, "SUBREDDIT")
	if err != nil {
		return nil, err
	}

	err = config.GetStringFromENV(&config.UserAgent, "USER_AGENT")
	if err != nil {
		return nil, err
	}

	err = config.GetStringFromENV(&config.Username, "USER_NAME")
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (cfg *Config) GetStringFromENV(field *string, variableName string) error {
	s := os.Getenv(variableName)
	if s == "" {
		return fmt.Errorf("%s was empty in env file", variableName)
	}
	*field = s
	return nil
}
