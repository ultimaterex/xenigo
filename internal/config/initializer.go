package config


import (
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

func EnsureConfigFile() {
    if _, err := os.Stat(defaultConfigFilename); os.IsNotExist(err) {
        defaultConfig := barebonesConfig
        err := os.WriteFile(defaultConfigFilename, defaultConfig, 0644)
        if err != nil {
            log.Fatalf("Error creating default config file: %v", err)
        }
    }
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
