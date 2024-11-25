// Package git provides git operations for version management.
// It handles git commands like tag creation, commit retrieval, and repository status checks.
package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

type git struct {
	tagPattern string
	workDir    string
}

// New creates a new git interface
func New(tagPattern, workDir string) Interface {
	return &git{
		tagPattern: tagPattern,
		workDir:    workDir,
	}
}

// GetLatestTag returns the latest tag that matches the pattern
func (g *git) GetLatestTag(pattern string) (string, error) {
	cmd := exec.Command("git", "tag", "--sort=-v:refname")
	cmd.Dir = g.workDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get tags: %v\n%s", err, stderr.String())
	}

	tags := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	if len(tags) == 0 || (len(tags) == 1 && tags[0] == "") {
		return "", fmt.Errorf("no tags found")
	}

	// If no pattern is provided, use the instance's pattern
	if pattern == "" {
		pattern = g.tagPattern
	}

	// Find the latest matching tag
	for _, tag := range tags {
		matched, err := filepath.Match(pattern, tag)
		if err != nil {
			return "", fmt.Errorf("invalid pattern: %v", err)
		}
		if matched {
			return tag, nil
		}
	}

	return "", fmt.Errorf("no matching tags found")
}

// GetCommitsSinceTag returns all commits since the given tag
func (g *git) GetCommitsSinceTag(tag string) ([]string, error) {
	cmd := exec.Command("git", "log", "--format=%B", tag+"..HEAD")
	cmd.Dir = g.workDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to get commits: %v\n%s", err, stderr.String())
	}

	commits := strings.Split(strings.TrimSpace(stdout.String()), "\n\n")
	var filtered []string
	for _, commit := range commits {
		if commit = strings.TrimSpace(commit); commit != "" {
			filtered = append(filtered, commit)
		}
	}
	return filtered, nil
}

// GetCommitsSinceTagForPath returns all commits since the given tag for the specified path
func (g *git) GetCommitsSinceTagForPath(tag, path string) ([]string, error) {
	cmd := exec.Command("git", "log", "--format=%B", tag+"..HEAD", "--", path)
	cmd.Dir = g.workDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to get commits: %v\n%s", err, stderr.String())
	}

	commits := strings.Split(strings.TrimSpace(stdout.String()), "\n\n")
	var filtered []string
	for _, commit := range commits {
		if commit = strings.TrimSpace(commit); commit != "" {
			filtered = append(filtered, commit)
		}
	}
	return filtered, nil
}

// GetFirstCommit returns the hash of the first commit
func (g *git) GetFirstCommit() (string, error) {
	cmd := exec.Command("git", "rev-list", "--max-parents=0", "HEAD")
	cmd.Dir = g.workDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get first commit: %v\n%s", err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}

// HasChanges returns true if there are uncommitted changes
func (g *git) HasChanges() (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = g.workDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return false, fmt.Errorf("failed to get status: %v\n%s", err, stderr.String())
	}

	return stdout.String() != "", nil
}

// CreateTag creates a new git tag
func (g *git) CreateTag(tag string, message string) error {
	cmd := exec.Command("git", "tag", "-a", tag, "-m", message)
	cmd.Dir = g.workDir
	return cmd.Run()
}

// PushTag pushes a tag to the remote repository
func (g *git) PushTag(tag string) error {
	cmd := exec.Command("git", "push", "origin", tag)
	cmd.Dir = g.workDir
	return cmd.Run()
}

// GetCurrentBranch returns the name of the current branch
func (g *git) GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = g.workDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get current branch: %v\n%s", err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}

// IsClean returns true if the repository has no uncommitted changes
func (g *git) IsClean() (bool, error) {
	hasChanges, err := g.HasChanges()
	if err != nil {
		return false, err
	}
	return !hasChanges, nil
}

// matchPattern checks if a string matches a glob pattern
func matchPattern(pattern, str string) (bool, error) {
	// For exact match
	if pattern == str {
		return true, nil
	}

	// For glob pattern match
	pattern = strings.TrimPrefix(pattern, "^")
	pattern = strings.TrimSuffix(pattern, "$")
	return strings.HasPrefix(str, pattern), nil
}

// GetCurrentVersion returns the current version based on the latest tag
func (g *git) GetCurrentVersion() (string, error) {
	return g.GetLatestTag(g.tagPattern)
}

// GetCommitsSinceVersion returns all commits since the given version
func (g *git) GetCommitsSinceVersion(version string) ([]string, error) {
	return g.GetCommitsSinceTag(version)
}
