// Package version provides version calculation and management functionality.
// It handles semantic versioning operations based on commit history and version formats.
package version

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/crazywolf132/bumpit/internal/config"
)

// Version handles version operations
type Version struct {
	cfg *config.Config
}

// New creates a new Version instance
func New(cfg *config.Config) *Version {
	return &Version{cfg: cfg}
}

// Calculate calculates the version based on current version (if any) and commits
func (v *Version) Calculate(currentVersion string, isInitial bool, commits []string) (string, error) {
	if isInitial || currentVersion == "" {
		version, err := CalculateInitialVersion(v.cfg, commits)
		if err != nil {
			return "", err
		}
		return version, nil
	}

	version, err := CalculateNextVersion(currentVersion, v.cfg, commits)
	if err != nil {
		return "", err
	}
	return version, nil
}

// IsValidVersion checks if a version string is valid
func (v *Version) IsValidVersion(version string) error {
	// Remove any path prefix
	versionNumber := version
	if strings.Contains(version, "/") {
		parts := strings.Split(version, "/")
		versionNumber = parts[len(parts)-1]
	}

	// Check if version starts with 'v'
	if !strings.HasPrefix(versionNumber, "v") {
		return fmt.Errorf("version must start with 'v'")
	}

	// Remove 'v' prefix for semver validation
	_, err := semver.NewVersion(strings.TrimPrefix(versionNumber, "v"))
	return err
}

// CompareVersions compares two version strings
func (v *Version) CompareVersions(v1, v2 string) (int, error) {
	ver1, err := semver.NewVersion(v1)
	if err != nil {
		return 0, fmt.Errorf("invalid version v1: %v", err)
	}

	ver2, err := semver.NewVersion(v2)
	if err != nil {
		return 0, fmt.Errorf("invalid version v2: %v", err)
	}

	if ver1.GreaterThan(ver2) {
		return 1, nil
	}
	if ver1.LessThan(ver2) {
		return -1, nil
	}
	return 0, nil
}

// CalculateInitialVersion calculates the initial version based on commit messages
func CalculateInitialVersion(cfg *config.Config, commits []string) (string, error) {
	major := 0
	minor := 0
	patch := 0

	for _, commit := range commits {
		if containsAny(commit, cfg.CommitTypes.Major) {
			major = 1
			minor = 0
			patch = 0
			break
		}
	}

	if major == 0 {
		for _, commit := range commits {
			if containsAny(commit, cfg.CommitTypes.Minor) {
				minor = 1
				patch = 0
				break
			}
		}
	}

	if major == 0 && minor == 0 {
		for _, commit := range commits {
			if containsAny(commit, cfg.CommitTypes.Patch) {
				patch = 1
				break
			}
		}
		// If no commit type matches, default to minor
		if patch == 0 {
			minor = 1
		}
	}

	return fmt.Sprintf("v%d.%d.%d", major, minor, patch), nil
}

// CalculateNextVersion calculates the next version based on the current version and commit messages
func CalculateNextVersion(currentVersion string, cfg *config.Config, commits []string) (string, error) {
	// Extract prefix if it exists
	prefix := ""
	versionNumber := currentVersion
	if strings.Contains(currentVersion, "/") {
		parts := strings.Split(currentVersion, "/")
		prefix = strings.Join(parts[:len(parts)-1], "/") + "/"
		versionNumber = parts[len(parts)-1]
	}

	// Remove 'v' prefix if it exists for semver parsing
	versionNumber = strings.TrimPrefix(versionNumber, "v")

	v, err := semver.NewVersion(versionNumber)
	if err != nil {
		return "", fmt.Errorf("invalid version format: %v", err)
	}

	major := v.Major()
	minor := v.Minor()
	patch := v.Patch()

	// Check for major version bump
	for _, commit := range commits {
		if containsAny(commit, cfg.CommitTypes.Major) {
			major++
			minor = 0
			patch = 0
			break
		}
	}

	// If no major bump, check for minor
	if major == v.Major() {
		for _, commit := range commits {
			if containsAny(commit, cfg.CommitTypes.Minor) {
				minor++
				patch = 0
				break
			}
		}
	}

	// If no major or minor bump, check for patch
	if major == v.Major() && minor == v.Minor() {
		for _, commit := range commits {
			if containsAny(commit, cfg.CommitTypes.Patch) {
				patch++
				break
			}
		}
		// If no commit type matches, default to patch
		if patch == v.Patch() && len(commits) > 0 {
			patch++
		}
	}

	return fmt.Sprintf("%sv%d.%d.%d", prefix, major, minor, patch), nil
}

// ValidateVersion validates if the version string is a valid semantic version
func ValidateVersion(version string) error {
	_, err := semver.NewVersion(version)
	return err
}

// CompareVersions compares two version strings and returns:
//
//	1 if v1 > v2
//	-1 if v1 < v2
//	0 if v1 == v2
func CompareVersions(v1, v2 string) (int, error) {
	ver1, err := semver.NewVersion(v1)
	if err != nil {
		return 0, fmt.Errorf("invalid version v1: %v", err)
	}

	ver2, err := semver.NewVersion(v2)
	if err != nil {
		return 0, fmt.Errorf("invalid version v2: %v", err)
	}

	if ver1.GreaterThan(ver2) {
		return 1, nil
	} else if ver1.LessThan(ver2) {
		return -1, nil
	}
	return 0, nil
}

// containsAny checks if any of the patterns are contained in the text
func containsAny(text string, patterns []string) bool {
	for _, pattern := range patterns {
		if strings.Contains(text, pattern) {
			return true
		}
	}
	return false
}
