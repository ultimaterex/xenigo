package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Developer Flags
const (
	DefaultSendFullConfigToLog    = false
	DefaultObfuscateConfigSecrets = true
	DefaultIgnoreCache            = false
	DefaultNotifyMute             = false
	DefaultForceSendInitial       = false
	DefaultDebugLogFileStructure  = false
)

type DeveloperFlags struct {
	SendFullConfigToLog    *bool `yaml:"send_full_config_to_log,omitempty"`
	ObfuscateConfigSecrets *bool `yaml:"obfuscate_config_secrets,omitempty"`
	IgnoreCache            *bool `yaml:"ignore_cache,omitempty"`
	NotifyMute             *bool `yaml:"notify_mute,omitempty"`
	ForceSendInitial       *bool `yaml:"force_send_initial,omitempty"`
	DebugLogFileStructure  *bool `yaml:"debug_log_file_structure,omitempty"`
}

func setDeveloperFlagsDefaults(config *Config) {
	if config.DeveloperFlags == nil {
		config.DeveloperFlags = &DeveloperFlags{
			SendFullConfigToLog:    boolPtr(DefaultSendFullConfigToLog),
			ObfuscateConfigSecrets: boolPtr(DefaultObfuscateConfigSecrets),
			IgnoreCache:            boolPtr(DefaultIgnoreCache),
			NotifyMute:             boolPtr(DefaultNotifyMute),
			ForceSendInitial:       boolPtr(DefaultForceSendInitial),
			DebugLogFileStructure:  boolPtr(DefaultDebugLogFileStructure),
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
	if isFlagSet(config.DeveloperFlags.DebugLogFileStructure) {
		log.Printf("Developer Flag Set: ForceSendInitial: %v", *config.DeveloperFlags.DebugLogFileStructure)
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
		config.DeveloperFlags.ForceSendInitial = boolPtr(DefaultForceSendInitial) // Default value for the new flag
	}
	if !isFlagSet(config.DeveloperFlags.DebugLogFileStructure) {
		config.DeveloperFlags.DebugLogFileStructure = boolPtr(DefaultDebugLogFileStructure) // Default value for the new flag
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

func logFileStructureBase(root string) {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		log.Printf("File: %s", path)
		return nil
	})
	if err != nil {
		log.Printf("Error walking the path %q: %v\n", root, err)
	}
}

func logFileStructure(root string, maxDepth int) {
    log.Printf("Logging file structure for root: %s up to depth: %d", root, maxDepth)
    logFileStructureHelper(root, 0, maxDepth)
}

func logFileStructureHelper(path string, currentDepth, maxDepth int) {
    if currentDepth > maxDepth {
        return
    }

    files, err := os.ReadDir(path)
    if err != nil {
        log.Printf("Error reading directory %q: %v\n", path, err)
        return
    }

    for _, file := range files {
        if file.IsDir() {
            log.Printf("%s[DIR] %s", strings.Repeat("  ", currentDepth), file.Name())
            logFileStructureHelper(filepath.Join(path, file.Name()), currentDepth+1, maxDepth)
        } else {
            log.Printf("%s[FILE] %s", strings.Repeat("  ", currentDepth), file.Name())
        }
    }
}

// Get flag value
func GetFlag(flag *bool) bool {
	return *flag
}

// / check if flag is not null
func isFlagSet(flag *bool) bool {
	return flag != nil
}
