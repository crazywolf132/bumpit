package git

// Interface defines the operations needed by bumpit
type Interface interface {
	GetLatestTag(pattern string) (string, error)
	GetCommitsSinceTag(tag string) ([]string, error)
	GetCommitsSinceTagForPath(tag, path string) ([]string, error)
	GetFirstCommit() (string, error)
	HasChanges() (bool, error)
	IsClean() (bool, error)
	CreateTag(tag string, message string) error
	PushTag(tag string) error
	GetCurrentBranch() (string, error)
	GetCurrentVersion() (string, error)
	GetCommitsSinceVersion(version string) ([]string, error)
}
