package version

import (
	"testing"

	"github.com/crazywolf132/bumpit/internal/config"
)

func newTestConfig() *config.Config {
	return &config.Config{
		VersionPrefix: "v",
		VersionFormat: "{major}.{minor}.{patch}",
		CommitTypes: config.CommitTypes{
			Major: []string{"BREAKING CHANGE"},
			Minor: []string{"feat"},
			Patch: []string{"fix"},
		},
	}
}

func TestCalculateInitialVersion(t *testing.T) {
	tests := []struct {
		name    string
		commits []string
		want    string
		wantErr bool
	}{
		{
			name:    "major change",
			commits: []string{"BREAKING CHANGE: something", "feat: new feature"},
			want:    "v1.0.0",
		},
		{
			name:    "minor change",
			commits: []string{"feat: new feature", "fix: bug fix"},
			want:    "v0.1.0",
		},
		{
			name:    "patch change",
			commits: []string{"fix: bug fix", "chore: update deps"},
			want:    "v0.0.1",
		},
		{
			name:    "no commits",
			commits: []string{},
			want:    "v0.1.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New(newTestConfig())
			got, err := v.Calculate("", true, tt.commits)
			if (err != nil) != tt.wantErr {
				t.Errorf("Calculate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Calculate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateNextVersion(t *testing.T) {
	tests := []struct {
		name    string
		current string
		commits []string
		want    string
		wantErr bool
	}{
		{
			name:    "major bump",
			current: "v1.0.0",
			commits: []string{"BREAKING CHANGE: something", "feat: new feature"},
			want:    "v2.0.0",
		},
		{
			name:    "minor bump",
			current: "v1.0.0",
			commits: []string{"feat: new feature", "fix: bug fix"},
			want:    "v1.1.0",
		},
		{
			name:    "patch bump",
			current: "v1.0.0",
			commits: []string{"fix: bug fix", "chore: update deps"},
			want:    "v1.0.1",
		},
		{
			name:    "no commits",
			current: "v1.0.0",
			commits: []string{},
			want:    "v1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New(newTestConfig())
			got, err := v.Calculate(tt.current, false, tt.commits)
			if (err != nil) != tt.wantErr {
				t.Errorf("Calculate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Calculate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		wantErr bool
	}{
		{
			name:    "valid version",
			version: "v1.0.0",
			wantErr: false,
		},
		{
			name:    "valid version with pre-release",
			version: "v1.0.0-alpha.1",
			wantErr: false,
		},
		{
			name:    "invalid version",
			version: "invalid",
			wantErr: true,
		},
		{
			name:    "missing v prefix",
			version: "1.0.0",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New(newTestConfig())
			err := v.IsValidVersion(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsValidVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name    string
		v1      string
		v2      string
		want    int
		wantErr bool
	}{
		{
			name: "v1 greater than v2",
			v1:   "v2.0.0",
			v2:   "v1.0.0",
			want: 1,
		},
		{
			name: "v1 less than v2",
			v1:   "v1.0.0",
			v2:   "v2.0.0",
			want: -1,
		},
		{
			name: "v1 equals v2",
			v1:   "v1.0.0",
			v2:   "v1.0.0",
			want: 0,
		},
		{
			name:    "invalid version v1",
			v1:      "invalid",
			v2:      "v1.0.0",
			wantErr: true,
		},
		{
			name:    "invalid version v2",
			v1:      "v1.0.0",
			v2:      "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New(newTestConfig())
			got, err := v.CompareVersions(tt.v1, tt.v2)
			if (err != nil) != tt.wantErr {
				t.Errorf("CompareVersions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got != tt.want {
				t.Errorf("CompareVersions() = %v, want %v", got, tt.want)
			}
		})
	}
}
