package main

import (
	"errors"
	"os"
	"regexp"

	"gopkg.in/yaml.v2"
)

type OAuthConfig struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
}

type FetchConfig struct {
	Subreddit string `yaml:"subreddit"`
	Sorting   string `yaml:"sorting"`
}

type Config struct {
	UserAgent string        `yaml:"user_agent"`
	OAuth     *OAuthConfig  `yaml:"oauth,omitempty"`
	Fetches   []FetchConfig `yaml:"fetches"`
}

func loadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Remove comments using a regular expression
	re := regexp.MustCompile(`(?m)^\s*#.*$|(?m)\s+#.*$`)
	cleanData := re.ReplaceAllString(string(data), "")

	var config Config
	if err := yaml.Unmarshal([]byte(cleanData), &config); err != nil {
		return nil, err
	}

	if config.OAuth != nil {
		if config.OAuth.ClientID == "" || config.OAuth.ClientSecret == "" || config.OAuth.Username == "" || config.OAuth.Password == "" {
			return nil, errors.New("oauth block is not correctly configured")
		}
	}

	return &config, nil
}

func determineContext(config *Config) string {
	if config.OAuth != nil {
		return "elevated"
	}
	return "standard"
}
