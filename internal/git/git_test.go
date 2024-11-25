package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func setupTestRepo(t *testing.T) (string, func()) {
	t.Helper()

	// Create a temporary directory for the test repository
	dir, err := os.MkdirTemp("", "git-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Initialize git repo
	cmds := [][]string{
		{"git", "init"},
		{"git", "config", "user.email", "test@example.com"},
		{"git", "config", "user.name", "Test User"},
		{"git", "config", "commit.gpgsign", "false"},
		{"git", "config", "tag.sort", "version:refname"},
	}

	for _, cmd := range cmds {
		c := exec.Command(cmd[0], cmd[1:]...)
		c.Dir = dir
		if err := c.Run(); err != nil {
			os.RemoveAll(dir)
			t.Fatalf("Failed to run command %v: %v", cmd, err)
		}
	}

	// Create initial commit
	testFile := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		os.RemoveAll(dir)
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd := exec.Command("git", "add", "test.txt")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		t.Fatalf("Failed to add test file: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		t.Fatalf("Failed to create initial commit: %v", err)
	}

	// Create test tags
	tags := []struct {
		name    string
		message string
	}{
		{"v1.0.0", "First release"},
		{"v1.1.0", "Second release"},
		{"v2.0.0", "Major release"},
		{"core/v1.0.0", "Core package release"},
	}

	for _, tag := range tags {
		cmd = exec.Command("git", "tag", "-a", tag.name, "-m", tag.message)
		cmd.Dir = dir
		if err := cmd.Run(); err != nil {
			os.RemoveAll(dir)
			t.Fatalf("Failed to create tag %s: %v", tag.name, err)
		}
	}

	cleanup := func() {
		os.RemoveAll(dir)
	}

	return dir, cleanup
}

func TestGetLatestTag(t *testing.T) {
	dir, cleanup := setupTestRepo(t)
	defer cleanup()

	g := &git{
		workDir: dir,
	}

	tests := []struct {
		name        string
		pattern     string
		want        string
		wantErr     bool
		errContains string
	}{
		{
			name:    "get latest version tag",
			pattern: "v*",
			want:    "v2.0.0",
		},
		{
			name:    "get latest core tag",
			pattern: "core/v*",
			want:    "core/v1.0.0",
		},
		{
			name:        "no matching tags",
			pattern:     "nonexistent*",
			wantErr:     true,
			errContains: "no matching tags found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := g.GetLatestTag(tt.pattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLatestTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				if tt.errContains != "" && err.Error() != tt.errContains {
					t.Errorf("GetLatestTag() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}
			if got != tt.want {
				t.Errorf("GetLatestTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCommitsSinceTag(t *testing.T) {
	dir, cleanup := setupTestRepo(t)
	defer cleanup()

	g := &git{
		workDir: dir,
	}

	// Add a new commit
	testFile := filepath.Join(dir, "test2.txt")
	if err := os.WriteFile(testFile, []byte("test2"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd := exec.Command("git", "add", "test2.txt")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to add test file: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", "feat: new feature")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create commit: %v", err)
	}

	commits, err := g.GetCommitsSinceTag("v2.0.0")
	if err != nil {
		t.Fatalf("GetCommitsSinceTag() error = %v", err)
	}

	if len(commits) != 1 {
		t.Errorf("GetCommitsSinceTag() = %v, want 1 commit", len(commits))
	}

	if commits[0] != "feat: new feature" {
		t.Errorf("GetCommitsSinceTag() = %v, want 'feat: new feature'", commits[0])
	}
}

func TestGetCommitsSinceTagForPath(t *testing.T) {
	dir, cleanup := setupTestRepo(t)
	defer cleanup()

	g := &git{
		workDir: dir,
	}

	// Create a subdirectory and add a commit
	subDir := filepath.Join(dir, "packages", "core")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	testFile := filepath.Join(subDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd := exec.Command("git", "add", "packages/core/test.txt")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to add test file: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", "feat(core): new core feature")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create commit: %v", err)
	}

	commits, err := g.GetCommitsSinceTagForPath("core/v1.0.0", "packages/core")
	if err != nil {
		t.Fatalf("GetCommitsSinceTagForPath() error = %v", err)
	}

	if len(commits) != 1 {
		t.Errorf("GetCommitsSinceTagForPath() = %v, want 1 commit", len(commits))
	}

	if commits[0] != "feat(core): new core feature" {
		t.Errorf("GetCommitsSinceTagForPath() = %v, want 'feat(core): new core feature'", commits[0])
	}
}
