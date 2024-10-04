package config

import (
	"encoding/json"
	"log"
)

// Developer Flags
const (
	DefaultSendFullConfigToLog    = false
	DefaultObfuscateConfigSecrets = true
	DefaultIgnoreCache            = false
	DefaultNotifyMute             = false
	DefaultForceSendInitial       = false
)

type DeveloperFlags struct {
	SendFullConfigToLog    *bool `yaml:"send_full_config_to_log,omitempty"`
	ObfuscateConfigSecrets *bool `yaml:"obfuscate_config_secrets,omitempty"`
	IgnoreCache            *bool `yaml:"ignore_cache,omitempty"`
	NotifyMute             *bool `yaml:"notify_mute,omitempty"`
	ForceSendInitial       *bool `yaml:"force_send_initial,omitempty"`
}

func setDeveloperFlagsDefaults(config *Config) {
	if config.DeveloperFlags == nil {
		config.DeveloperFlags = &DeveloperFlags{
			SendFullConfigToLog:    boolPtr(DefaultSendFullConfigToLog),
			ObfuscateConfigSecrets: boolPtr(DefaultObfuscateConfigSecrets),
			IgnoreCache:            boolPtr(DefaultIgnoreCache),
			NotifyMute:             boolPtr(DefaultNotifyMute),
			ForceSendInitial:       boolPtr(DefaultForceSendInitial),
		}
		return
	}

	// Log only the set developer flags
	if isFlagSet(config.DeveloperFlags.SendFullConfigToLog) {
		log.Printf("Developer Flag Set: SendFullConfigToLog: %v", *config.DeveloperFlags.SendFullConfigToLog)
	}
	if isFlagSet(config.DeveloperFlags.ObfuscateConfigSecrets) {
		log.Printf("Developer Flag Set: ObfuscateConfigSecrets: %v", *config.DeveloperFlags.ObfuscateConfigSecrets)
	}
	if isFlagSet(config.DeveloperFlags.IgnoreCache) {
		log.Printf("Developer Flag Set: IgnoreCache: %v", *config.DeveloperFlags.IgnoreCache)
	}
	if isFlagSet(config.DeveloperFlags.NotifyMute) {
		log.Printf("Developer Flag Set: NotifyMute: %v", *config.DeveloperFlags.NotifyMute)
	}
	if isFlagSet(config.DeveloperFlags.ForceSendInitial) {
		log.Printf("Developer Flag Set: ForceSendInitial: %v", *config.DeveloperFlags.ForceSendInitial)
	}

	// Set the default values for the developer flags
	if !isFlagSet(config.DeveloperFlags.SendFullConfigToLog) {
		config.DeveloperFlags.SendFullConfigToLog = boolPtr(DefaultSendFullConfigToLog)
	}
	if !isFlagSet(config.DeveloperFlags.ObfuscateConfigSecrets) {
		config.DeveloperFlags.ObfuscateConfigSecrets = boolPtr(DefaultObfuscateConfigSecrets)
	}
	if !isFlagSet(config.DeveloperFlags.IgnoreCache) {
		config.DeveloperFlags.IgnoreCache = boolPtr(DefaultIgnoreCache)
	}
	if !isFlagSet(config.DeveloperFlags.NotifyMute) {
		config.DeveloperFlags.NotifyMute = boolPtr(DefaultNotifyMute)
	}
	if !isFlagSet(config.DeveloperFlags.ForceSendInitial) {
		config.DeveloperFlags.ForceSendInitial = boolPtr(false) // Default value for the new flag
	}
}

func obfuscateSecrets(config *Config) {
	if config.OAuth != nil {
		config.OAuth.ClientID = "********"
		config.OAuth.ClientSecret = "********"
		config.OAuth.Username = ""
		config.OAuth.Password = "********"
	}
	for i := range config.Targets {
		config.Targets[i].Output.WebhookURL = "********"
	}
}

func logFullConfig(config *Config) {
	configCopy := *config

	if GetFlag(config.DeveloperFlags.ObfuscateConfigSecrets) {
		obfuscateSecrets(&configCopy)
	}

	configJSON, err := json.MarshalIndent(configCopy, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal config to JSON: %v", err)
		return
	}
	log.Printf("Full Config:\n%s", string(configJSON))
}


// Get flag value
func GetFlag(flag *bool) bool {
	return *flag
}

/// check if flag is not null
func isFlagSet(flag *bool) bool {
	return flag != nil
}
