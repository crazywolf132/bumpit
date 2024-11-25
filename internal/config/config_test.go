package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "bumpit-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test config files
	basicConfig := `
version_prefix: "v"
version_format: "{major}.{minor}.{patch}"
pre_release: ""
build_metadata: ""
default_command: "patch"
commit_types:
  major:
    - "BREAKING CHANGE"
  minor:
    - "feat"
  patch:
    - "fix"
git:
  tag_pattern: "v*"
  auto_push: false
output:
  debug: false
  color: true
paths:
  - path: "core"
    version_prefix: "core/v"
    version_format: "{major}.{minor}.{patch}"
    pre_release: ""
    build_metadata: ""
    tag_pattern: "core/v*"
    default_command: "patch"
    commit_types:
      major:
        - "BREAKING CHANGE"
      minor:
        - "feat"
      patch:
        - "fix"
`

	configPath := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(basicConfig), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	tests := []struct {
		name    string
		envVar  string
		wantErr bool
	}{
		{
			name:    "basic config",
			envVar:  configPath,
			wantErr: false,
		},
		{
			name:    "core package config",
			envVar:  configPath,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("BUMPIT_CONFIG", tt.envVar)
			defer os.Unsetenv("BUMPIT_CONFIG")

			got, err := LoadConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if got.VersionFormat != "{major}.{minor}.{patch}" {
					t.Errorf("LoadConfig() version format = %v, want {major}.{minor}.{patch}", got.VersionFormat)
				}
				if got.VersionPrefix != "v" {
					t.Errorf("LoadConfig() version prefix = %v, want v", got.VersionPrefix)
				}
			}
		})
	}
}

func TestConfigValidation(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "bumpit-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	validConfig := `
version_prefix: "v"
version_format: "{major}.{minor}.{patch}"
paths:
  - path: "core"
    version_prefix: "core/v"
    version_format: "{major}.{minor}.{patch}"
`

	invalidConfig := `
version_prefix: "v"
version_format: "x.y.z"
`

	tests := []struct {
		name     string
		config   string
		wantErr  bool
		errCheck func(error) bool
	}{
		{
			name:    "valid config",
			config:  validConfig,
			wantErr: false,
		},
		{
			name:    "valid path config",
			config:  validConfig,
			wantErr: false,
		},
		{
			name:    "invalid version format",
			config:  invalidConfig,
			wantErr: true,
			errCheck: func(err error) bool {
				return err != nil && err.Error() == "invalid version format: must contain {major}, {minor}, and {patch}"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configPath := filepath.Join(tmpDir, "config.yaml")
			if err := os.WriteFile(configPath, []byte(tt.config), 0644); err != nil {
				t.Fatalf("Failed to write config file: %v", err)
			}

			os.Setenv("BUMPIT_CONFIG", configPath)
			defer os.Unsetenv("BUMPIT_CONFIG")

			_, err := LoadConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.errCheck != nil && !tt.errCheck(err) {
				t.Errorf("LoadConfig() error = %v, does not match expected error condition", err)
			}
		})
	}
}
