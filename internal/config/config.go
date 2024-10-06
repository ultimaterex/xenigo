package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v2"
)

// Options
const (
	DefaultInterval      = 60
	DefaultLimit         = 3
	DefaultRetryCount    = 3
	DefaultRetryInterval = 2
)



type OAuthConfig struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
}

type Options struct {
	Interval       int  `yaml:"interval"`
	Limit          int  `yaml:"limit"`
	RetryCount     int  `yaml:"retry_count"`
	RetryInterval  int  `yaml:"retry_interval"`
	EnableFallback bool `yaml:"enable_fallback"`
}

type Config struct {
	UserAgent      string          `yaml:"user_agent"`
	OAuth          *OAuthConfig    `yaml:"oauth,omitempty"`
	Targets        []Target        `yaml:"targets"`
	Options        *Options        `yaml:"options,omitempty"`
	DeveloperFlags *DeveloperFlags `yaml:"developer_flags,omitempty"`
}

type Context string

const (
	ContextStandard Context = "standard"
	ContextElevated Context = "elevated"
)

type AppConfig struct {
	Config  *Config
	Context Context
}

type Target struct {
	Name    string `yaml:"name"`
	Monitor struct {
		Subreddit string `yaml:"subreddit"`
		Sorting   string `yaml:"sorting"`
	} `yaml:"monitor"`
	Output  OutputConfig `yaml:"output"`
	Options *Options     `yaml:"options,omitempty"`
}

type OutputType string

const (
	OutputTypeDiscord OutputType = "discord"
	OutputTypeSlack   OutputType = "slack"
)

type OutputConfig struct {
	Type       OutputType `yaml:"type"`
	WebhookURL string     `yaml:"webhook_url"`
	Format     struct {
		URL           *bool `yaml:"url"`
		Author        *bool `yaml:"author"`
		Subreddit     *bool `yaml:"subreddit"`
		DiscussionURL *bool `yaml:"discussion_url"`
	} `yaml:"format"`
}

func loadConfigFile(filename string) (*Config, error) {
	possiblePaths := []string{
        filename,
        filepath.Join("config", filename),
        filepath.Join("data", filename),
    }

    var data []byte
    var err error
    for _, path := range possiblePaths {
        data, err = os.ReadFile(path)
        if err == nil {
            break
        }
    }

    if err != nil {
		// if GetFlag(config.DeveloperFlags.DebugLogFileStructure) { #TODO config isn't initialized yet...
		if true {
            logFileStructure(".")
        }
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

	// Remove yaml comments using a regular expression
	re := regexp.MustCompile(`(?m)^\s*#.*$|(?m)\s+#.*$`)
	cleanData := re.ReplaceAllString(string(data), "")

	var config Config
	if err := yaml.Unmarshal([]byte(cleanData), &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	setGlobalDefaults(&config)
	setDeveloperFlagsDefaults(&config)

	for i, target := range config.Targets {
		if target.Monitor.Subreddit == "" || target.Monitor.Sorting == "" {
			return nil, errors.New("monitor block is not correctly configured")
		}
		if target.Output.WebhookURL == "" {
			return nil, errors.New("output block is not correctly configured")
		}
		
		if target.Name == "" {
			config.Targets[i].Name = target.Monitor.Subreddit
		}
		// Initialize Format using the helper function
		config.Targets[i].Output.Format = initializeFormat(target.Output.Format)
		// Check if all format options are set to false
		if !*config.Targets[i].Output.Format.Subreddit && !*config.Targets[i].Output.Format.Author && !*config.Targets[i].Output.Format.DiscussionURL && !*config.Targets[i].Output.Format.URL {
			return nil, fmt.Errorf("all format options are set to false for target %s, which will cause problems", target.Name)
		}
		setTargetDefaults(&config, &config.Targets[i])
	}

	if GetFlag(config.DeveloperFlags.SendFullConfigToLog) {
		logFullConfig(&config)
	}

	return &config, nil
}



func validateConfig(config *Config) error {
	if config.UserAgent == "" {
		return errors.New("user_agent is required")
	}
	if config.OAuth != nil {
		if config.OAuth.ClientID == "" || config.OAuth.ClientSecret == "" || config.OAuth.Username == "" || config.OAuth.Password == "" {
			return errors.New("oauth block is not correctly configured")
		}
	}
	for _, target := range config.Targets {
		if target.Monitor.Subreddit == "" || target.Monitor.Sorting == "" {
			return errors.New("monitor block is not correctly configured")
		}
		if target.Output.WebhookURL == "" {
			return errors.New("output block is not correctly configured")
		}
	}
	return nil
}

func setGlobalDefaults(config *Config) {
	if config.Options == nil {
		config.Options = &Options{
			Interval:      DefaultInterval,
			Limit:         DefaultLimit,
			RetryCount:    DefaultRetryCount,
			RetryInterval: DefaultRetryInterval,
		}
	}
}



func setTargetDefaults(config *Config, target *Target) {
	if target.Options == nil {
		target.Options = &Options{}
	}
	if config.Options != nil {
		if target.Options.Interval == 0 {
			target.Options.Interval = config.Options.Interval
		}
		if target.Options.Limit == 0 {
			target.Options.Limit = config.Options.Limit
		}
		if target.Options.RetryCount == 0 {
			target.Options.RetryCount = config.Options.RetryCount
		}
		if target.Options.RetryInterval == 0 {
			target.Options.RetryInterval = config.Options.RetryInterval
		}
	}
}

func boolPtr(b bool) *bool {
	return &b
}

