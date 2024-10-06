package config

import (
    "fmt"
    "log"
    "os"
)

const defaultConfigFilename = "config.yaml"

var barebonesConfig = []byte(`
user_agent: xenigo
oauth:
  client_id: your_client_id
  client_secret: your_client_secret
  username: your_username
  password: your_password
options:
  interval: 60
  limit: 3
  retry_count: 5
  retry_interval: 5
targets:
  - monitor:
      subreddit: hardwareswap
      sorting: new
     output:
      type: discord
      webhook_url: your_webhook_url
      format:
       subreddit: false
        author: false
        discussion_url: false
`)

func EnsureConfigFile(filename string) error {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return createDefaultConfigFile(filename)
    }
    if err != nil {
        return fmt.Errorf("error checking config file: %w", err)
    }
    if info.IsDir() {
        return fmt.Errorf("config file is a directory")
    }
    file, err := os.Open(filename)
    if err != nil {
        return fmt.Errorf("error opening config file: %w", err)
    }
    defer file.Close()
    return nil
}

func createDefaultConfigFile(filename string) error {
    err := os.WriteFile(filename, barebonesConfig, 0644)
    if err != nil {
        return fmt.Errorf("error creating default config file: %w", err)
    }
    log.Printf("Created default config file: %s", filename)
    return nil
}

func initializeFormat(format struct {
	URL           *bool `yaml:"url"`
	Author        *bool `yaml:"author"`
	Subreddit     *bool `yaml:"subreddit"`
	DiscussionURL *bool `yaml:"discussion_url"`
}) struct {
	URL           *bool `yaml:"url"`
	Author        *bool `yaml:"author"`
	Subreddit     *bool `yaml:"subreddit"`
	DiscussionURL *bool `yaml:"discussion_url"`
} {
	if format.URL == nil {
		format.URL = boolPtr(true)
	}
	if format.Author == nil {
		format.Author = boolPtr(true)
	}
	if format.Subreddit == nil {
		format.Subreddit = boolPtr(true)
	}
	if format.DiscussionURL == nil {
		format.DiscussionURL = boolPtr(true)
	}
	return format
}

func LoadConfig() (*AppConfig, error) {
  config, err := loadConfigFile(defaultConfigFilename)
  if err != nil {
      return nil, err
  }

	context := ContextStandard
	if config.OAuth != nil {
		context = ContextElevated
	}

	log.Printf("Running in context: %s", context)

	return &AppConfig{Config: config, Context: context}, nil
}
