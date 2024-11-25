// Package config provides configuration management for the bumpit tool.
// It handles loading and validating configuration from various sources.
package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

// Config represents the main configuration structure for bumpit.
type Config struct {
	VersionPrefix  string       `yaml:"version_prefix"`
	VersionFormat  string       `yaml:"version_format"`
	PreRelease     string       `yaml:"pre_release"`
	BuildMetadata  string       `yaml:"build_metadata"`
	DefaultCommand string       `yaml:"default_command"`
	CommitTypes    CommitTypes  `yaml:"commit_types"`
	Git            GitConfig    `yaml:"git"`
	Output         OutputConfig `yaml:"output"`
	Paths          []PathConfig `yaml:"paths"`
}

// CommitTypes defines which commit message prefixes trigger different types of version bumps.
type CommitTypes struct {
	Major []string `yaml:"major"`
	Minor []string `yaml:"minor"`
	Patch []string `yaml:"patch"`
}

// GitConfig holds git-specific configuration options.
type GitConfig struct {
	TagPattern string `yaml:"tag_pattern"`
	AutoPush   bool   `yaml:"auto_push"`
}

// OutputConfig defines output formatting options.
type OutputConfig struct {
	Debug bool `yaml:"debug"`
	Color bool `yaml:"color"`
}

// PathConfig represents configuration for a specific path in the repository.
type PathConfig struct {
	Path           string      `yaml:"path"`
	VersionPrefix  string      `yaml:"version_prefix"`
	VersionFormat  string      `yaml:"version_format"`
	PreRelease     string      `yaml:"pre_release"`
	BuildMetadata  string      `yaml:"build_metadata"`
	TagPattern     string      `yaml:"tag_pattern"`
	DefaultCommand string      `yaml:"default_command"`
	CommitTypes    CommitTypes `yaml:"commit_types"`
}

// LoadConfig loads the configuration from various sources and validates it.
// It first tries to load from the BUMPIT_CONFIG environment variable,
// then from default locations, and finally applies default values.
func LoadConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigType("yaml")

	// Set defaults
	v.SetDefault("version_format", "{major}.{minor}.{patch}")
	v.SetDefault("version_prefix", "v")
	v.SetDefault("commit_types", map[string][]string{
		"major": {"BREAKING CHANGE"},
		"minor": {"feat"},
		"patch": {"fix"},
	})

	// Check for config file path in environment variable
	configPath := os.Getenv("BUMPIT_CONFIG")
	if configPath != "" {
		// Read from specified file
		configData, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %v", err)
		}
		if err := v.ReadConfig(strings.NewReader(string(configData))); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %v", err)
		}
	} else {
		// Load from default locations
		v.SetConfigName("default.config")
		v.AddConfigPath(".")
		v.AddConfigPath("$HOME/.config/bumpit")

		// Load default config
		if err := v.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return nil, err
			}
		}

		// Look for and merge with .bumpit.yaml in the current directory
		v.SetConfigName(".bumpit")
		if err := v.MergeInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return nil, err
			}
		}
	}

	// Check raw version format before unmarshaling
	if v.IsSet("version_format") {
		rawFormat := v.GetString("version_format")
		if err := validateVersionFormat(rawFormat); err != nil {
			return nil, err
		}
	}

	// Check raw version format for paths before unmarshaling
	if v.IsSet("paths") {
		paths := v.Get("paths").([]interface{})
		for _, path := range paths {
			pathMap := path.(map[string]interface{})
			if format, ok := pathMap["version_format"]; ok {
				if err := validateVersionFormat(format.(string)); err != nil {
					return nil, fmt.Errorf("invalid version format for path %s: %v", pathMap["path"], err)
				}
			}
		}
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	// Set default values
	if config.VersionPrefix == "" {
		config.VersionPrefix = "v"
	}
	if config.VersionFormat == "" {
		config.VersionFormat = "{major}.{minor}.{patch}"
	}

	if len(config.CommitTypes.Major) == 0 {
		config.CommitTypes.Major = []string{"BREAKING CHANGE"}
	}
	if len(config.CommitTypes.Minor) == 0 {
		config.CommitTypes.Minor = []string{"feat"}
	}
	if len(config.CommitTypes.Patch) == 0 {
		config.CommitTypes.Patch = []string{"fix"}
	}

	// Set default values for paths
	for i := range config.Paths {
		if config.Paths[i].VersionPrefix == "" {
			config.Paths[i].VersionPrefix = config.VersionPrefix
		}
		if config.Paths[i].VersionFormat == "" {
			config.Paths[i].VersionFormat = config.VersionFormat
		}
		if len(config.Paths[i].CommitTypes.Major) == 0 {
			config.Paths[i].CommitTypes.Major = config.CommitTypes.Major
		}
		if len(config.Paths[i].CommitTypes.Minor) == 0 {
			config.Paths[i].CommitTypes.Minor = config.CommitTypes.Minor
		}
		if len(config.Paths[i].CommitTypes.Patch) == 0 {
			config.Paths[i].CommitTypes.Patch = config.CommitTypes.Patch
		}
	}

	return &config, nil
}

// GetCommitType determines the type of version bump needed based on commit message
func (c *Config) GetCommitType(commitMsg string) string {
	// Check major changes
	for _, prefix := range c.CommitTypes.Major {
		if hasPrefix(commitMsg, prefix) {
			return "major"
		}
	}

	// Check minor changes
	for _, prefix := range c.CommitTypes.Minor {
		if hasPrefix(commitMsg, prefix) {
			return "minor"
		}
	}

	// Check patch changes
	for _, prefix := range c.CommitTypes.Patch {
		if hasPrefix(commitMsg, prefix) {
			return "patch"
		}
	}

	return "none"
}

// hasPrefix checks if a message has a specific prefix
func hasPrefix(msg, prefix string) bool {
	return strings.Contains(strings.ToLower(msg), strings.ToLower(prefix))
}

// GetPathConfig returns the configuration for a specific path
func (c *Config) GetPathConfig(path string) PathConfig {
	// If no path is provided, return a default config
	if path == "" {
		return PathConfig{
			VersionPrefix: c.VersionPrefix,
			VersionFormat: c.VersionFormat,
			CommitTypes:   c.CommitTypes,
		}
	}

	// Find the most specific path configuration
	var bestMatch PathConfig
	bestMatchLen := -1

	for _, pathConfig := range c.Paths {
		// Check if the current path is under this configuration path
		rel, err := filepath.Rel(pathConfig.Path, path)
		if err != nil || strings.HasPrefix(rel, "..") {
			continue
		}

		// Get the length of the match (shorter relative path means better match)
		matchLen := len(rel)
		if matchLen < bestMatchLen || bestMatchLen == -1 {
			bestMatch = pathConfig
			bestMatchLen = matchLen
		}
	}

	// If no matching path found, return default config
	if bestMatchLen == -1 {
		return PathConfig{
			VersionPrefix: c.VersionPrefix,
			VersionFormat: c.VersionFormat,
			CommitTypes:   c.CommitTypes,
		}
	}

	// Fill in any missing values with defaults from the main config
	if bestMatch.VersionPrefix == "" {
		bestMatch.VersionPrefix = c.VersionPrefix
	}
	if bestMatch.VersionFormat == "" {
		bestMatch.VersionFormat = c.VersionFormat
	}
	if len(bestMatch.CommitTypes.Major) == 0 {
		bestMatch.CommitTypes.Major = c.CommitTypes.Major
	}
	if len(bestMatch.CommitTypes.Minor) == 0 {
		bestMatch.CommitTypes.Minor = c.CommitTypes.Minor
	}
	if len(bestMatch.CommitTypes.Patch) == 0 {
		bestMatch.CommitTypes.Patch = c.CommitTypes.Patch
	}

	return bestMatch
}

func validateVersionFormat(format string) error {
	if format == "" {
		return nil // Empty format will be replaced with default
	}
	if format == "x.y.z" {
		return fmt.Errorf("invalid version format: must contain {major}, {minor}, and {patch}")
	}
	if !strings.Contains(format, "{major}") ||
		!strings.Contains(format, "{minor}") ||
		!strings.Contains(format, "{patch}") {
		return fmt.Errorf("invalid version format: must contain {major}, {minor}, and {patch}")
	}
	return nil
}
