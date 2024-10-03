package config

import (
	"errors"
	"regexp"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name        string
		configData  string
		expectError bool
	}{
		{
			name: "Valid config with OAuth",
			configData: `
user_agent: xenigo
oauth:
  client_id: test_client_id
  client_secret: test_client_secret
  username: test_username
  password: test_password
fetches:
  - subreddit: buildapcsales
    sorting: new
  - subreddit: hardwareswap
    sorting: new
`,
			expectError: false,
		},
		{
			name: "Valid config without OAuth",
			configData: `
user_agent: xenigo
fetches:
  - subreddit: buildapcsales
    sorting: new
  - subreddit: hardwareswap
    sorting: new
`,
			expectError: false,
		},
		{
			name: "Invalid config with incomplete OAuth",
			configData: `
user_agent: xenigo
oauth:
  client_id: test_client_id
  client_secret: test_client_secret
fetches:
  - subreddit: buildapcsales
    sorting: new
  - subreddit: hardwareswap
    sorting: new
`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := loadConfigFromString(tt.configData)
			if (err != nil) != tt.expectError {
				t.Errorf("loadConfig() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if config == nil && !tt.expectError {
				t.Errorf("Expected valid config, got nil")
			}
		})
	}
}

func loadConfigFromString(data string) (*Config, error) {
	// Remove comments using a regular expression
	re := regexp.MustCompile(`(?m)^\s*#.*$|(?m)\s+#.*$`)
	cleanData := re.ReplaceAllString(data, "")

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
