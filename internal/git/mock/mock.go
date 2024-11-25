package mock

import (
	"github.com/crazywolf132/bumpit/internal/git"
)

// Git is a mock implementation of the git.Interface
type Git struct {
	// Mock data
	LastTag          string
	LastTagExists    bool
	LastTagError     error
	CommitsSinceTag  []string
	CommitsError     error
	FirstCommit      string
	FirstCommitError error
	HasChangesResult bool
	HasChangesError  error
	CurrentBranch    string
	BranchError      error
	IsCleanResult    bool
	IsCleanError     error
	CreateTagError   error
	PushTagError     error
	LatestTag        string
	LatestTagError   error
	CurrentVersion   string
	VersionError     error
	LatestTagFunc    func(pattern string) (string, error)
	CommitsSinceTagForPathFunc func(tag string, path string) ([]string, error)
}

// New creates a new mock Git instance
func New() git.Interface {
	return &Git{}
}

// GetLastTag returns mock data for the last tag
func (g *Git) GetLastTag() (string, bool, error) {
	return g.LastTag, g.LastTagExists, g.LastTagError
}

// GetCommitsSinceTag returns mock data for commits since tag
func (g *Git) GetCommitsSinceTag(tag string) ([]string, error) {
	return g.CommitsSinceTag, g.CommitsError
}

// GetCommitsSinceTagForPath returns commits since tag for a specific path
func (g *Git) GetCommitsSinceTagForPath(tag string, path string) ([]string, error) {
	if g.CommitsSinceTagForPathFunc != nil {
		return g.CommitsSinceTagForPathFunc(tag, path)
	}
	return []string{}, nil
}

// GetFirstCommit returns mock data for the first commit
func (g *Git) GetFirstCommit() (string, error) {
	return g.FirstCommit, g.FirstCommitError
}

// HasChanges returns mock data for uncommitted changes
func (g *Git) HasChanges() (bool, error) {
	return g.HasChangesResult, g.HasChangesError
}

// CreateTag returns mock data for tag creation
func (g *Git) CreateTag(tag string, message string) error {
	return g.CreateTagError
}

// PushTag returns mock data for tag pushing
func (g *Git) PushTag(tag string) error {
	return g.PushTagError
}

// GetCurrentBranch returns mock data for current branch
func (g *Git) GetCurrentBranch() (string, error) {
	return g.CurrentBranch, g.BranchError
}

// IsClean returns mock data for repository cleanliness
func (g *Git) IsClean() (bool, error) {
	return g.IsCleanResult, g.IsCleanError
}

// GetLatestTag returns mock data for the latest tag
func (g *Git) GetLatestTag(pattern string) (string, error) {
	if g.LatestTagFunc != nil {
		return g.LatestTagFunc(pattern)
	}
	return g.LatestTag, g.LatestTagError
}

// GetCurrentVersion returns mock data for current version
func (g *Git) GetCurrentVersion() (string, error) {
	return g.CurrentVersion, g.VersionError
}

// GetCommitsSinceVersion returns mock data for commits since version
func (g *Git) GetCommitsSinceVersion(version string) ([]string, error) {
	return g.GetCommitsSinceTag(version)
}
