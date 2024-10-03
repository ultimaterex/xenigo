package config

import (
    "errors"
    "fmt"
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

type Monitor struct {
    Subreddit string `yaml:"subreddit"`
    Sorting   string `yaml:"sorting"`
}

type Output struct {
    WebhookType string `yaml:"webhook_type"`
    WebhookURL  string `yaml:"webhook_url"`
}

type Options struct {
    Interval int `yaml:"interval"`
    Limit    int `yaml:"limit"`
}

type Target struct {
    Name    string  `yaml:"name,omitempty"`
    Monitor Monitor `yaml:"monitor"`
    Output  Output  `yaml:"output"`
    Options Options `yaml:"options"`
}

type Config struct {
    UserAgent string   `yaml:"user_agent"`
    OAuth     *OAuthConfig `yaml:"oauth,omitempty"`
    Targets   []Target `yaml:"targets"`
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

func loadConfigFile(filename string) (*Config, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    // Remove yaml comments using a regular expression
    re := regexp.MustCompile(`(?m)^\s*#.*$|(?m)\s+#.*$`)
    cleanData := re.ReplaceAllString(string(data), "")

    var config Config
    if err := yaml.Unmarshal([]byte(cleanData), &config); err != nil {
        return nil, fmt.Errorf("failed to unmarshal config: %w", err)
    }

    if config.OAuth != nil {
        if config.OAuth.ClientID == "" || config.OAuth.ClientSecret == "" || config.OAuth.Username == "" || config.OAuth.Password == "" {
            return nil, errors.New("oauth block is not correctly configured")
        }
    }

    for i, target := range config.Targets {
        if target.Monitor.Subreddit == "" || target.Monitor.Sorting == "" {
            return nil, errors.New("monitor block is not correctly configured")
        }
        if target.Output.WebhookType == "" || target.Output.WebhookURL == "" {
            return nil, errors.New("output block is not correctly configured")
        }
        if target.Name == "" {
            config.Targets[i].Name = target.Monitor.Subreddit
        }
    }

    return &config, nil
}

func LoadConfig() (*AppConfig, error) {
    filenames := []string{"config.yml", "config.yaml"}
    var config *Config
    var err error

    for _, filename := range filenames {
        config, err = loadConfigFile(filename)
        if err == nil {
            break
        }
    }

    if err != nil {
        return nil, err
    }

    context := ContextStandard
    if config.OAuth != nil {
        context = ContextElevated
    }

    return &AppConfig{Config: config, Context: context}, nil
}
