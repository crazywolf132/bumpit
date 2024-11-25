# Bumpit

> 🚀 Powerful, flexible, and language-agnostic semantic versioning automation for your projects

[![GitHub release (latest by date)](https://img.shields.io/github/v/release/crazywolf132/bumpit)](https://github.com/crazywolf132/bumpit/releases)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/crazywolf132/bumpit/ci.yml?branch=main)](https://github.com/crazywolf132/bumpit/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/crazywolf132/bumpit)](https://goreportcard.com/report/github.com/crazywolf132/bumpit)
[![codecov](https://codecov.io/gh/crazywolf132/bumpit/branch/main/graph/badge.svg)](https://codecov.io/gh/crazywolf132/bumpit)
[![Go Reference](https://pkg.go.dev/badge/github.com/crazywolf132/bumpit.svg)](https://pkg.go.dev/github.com/crazywolf132/bumpit)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## The Problem

Managing version numbers across projects can be tedious and error-prone:
- Manual version bumping leads to mistakes
- Different team members might bump versions inconsistently
- Conventional commits aren't being leveraged effectively
- Existing tools are often language-specific or too complex
- Custom scripts are hard to maintain and lack features

## The Solution

Bumpit is a lightweight, powerful tool that automates semantic versioning based on your git commits. It:
- 🎯 Automatically determines the next version based on conventional commits
- 🔧 Works with any programming language or project type
- ⚡️ Provides both a CLI tool and a GitHub Action
- 🎨 Offers flexible version formatting options
- 🔌 Executes custom commands with the new version
- 🤝 Integrates seamlessly with your existing workflow

### Why Bumpit?

#### vs semantic-release
- Lighter weight and faster
- Language-agnostic by design
- More flexible version formatting
- Simpler configuration
- Works without Node.js
- Custom command execution

#### vs custom scripts
- Battle-tested and maintained
- Proper conventional commit parsing
- Rich configuration options
- GitHub Action support
- No reinventing the wheel
- Zero dependencies in your project

## Quick Start

### CLI Installation
```bash
go install github.com/crazywolf132/bumpit@latest
```

### Basic Usage

1. **CLI**: Run bumpit with a command template
```bash
# Create a git tag
bumpit "git tag ${version}"

# Update package.json
bumpit "npm version ${version}"

# Update any version file
bumpit "sed -i '' 's/version = .*/version = \"${version}\"/' version.txt"
```

2. **GitHub Action**: Add to your workflow
```yaml
- name: Update version
  uses: crazywolf132/bumpit@v1
  with:
    github_token: ${{ secrets.GITHUB_TOKEN }}
```

## Configuration Deep Dive

### Version Format Control

Bumpit gives you complete control over your version format through the `.bumpit.yaml` configuration:

```yaml
# Version Formatting
version_prefix: "v"               # Prefix for version (e.g., v1.0.0)
version_format: "{major}.{minor}.{patch}"  # Version number format
pre_release: "beta.1"            # Pre-release suffix (e.g., v1.0.0-beta.1)
build_metadata: "20230815"       # Build metadata (e.g., v1.0.0+20230815)

# Commit Analysis
commit_types:
  major:                         # Triggers major version bump (1.0.0 -> 2.0.0)
    - "BREAKING CHANGE"
    - "major"
  minor:                         # Triggers minor version bump (1.0.0 -> 1.1.0)
    - "feat"
  patch:                         # Triggers patch version bump (1.0.0 -> 1.0.1)
    - "fix"
    - "chore"
    - "docs"
    - "style"
    - "refactor"
    - "perf"
    - "test"

# Behavior
default_command: "git tag ${version}"  # Default command if none specified
git:
  tag_pattern: "v*"              # Pattern for finding version tags
  auto_push: false              # Auto-push new tags

# Output
output:
  debug: false                  # Enable debug logging
  color: true                   # Colorize output
```

### Version Format Examples

```yaml
# Standard semver: v1.2.3
version_prefix: "v"
version_format: "{major}.{minor}.{patch}"

# No prefix: 1.2.3
version_prefix: ""
version_format: "{major}.{minor}.{patch}"

# Custom format: 0.1.2
version_prefix: ""
version_format: "0.{major}.{minor}"

# With pre-release: v1.2.3-alpha.1
version_prefix: "v"
version_format: "{major}.{minor}.{patch}"
pre_release: "alpha.1"

# With build metadata: v1.2.3+20230815
version_prefix: "v"
version_format: "{major}.{minor}.{patch}"
build_metadata: "20230815"

# Combined: v1.2.3-beta.1+20230815
version_prefix: "v"
version_format: "{major}.{minor}.{patch}"
pre_release: "beta.1"
build_metadata: "20230815"
```

## Advanced Use Cases

### Monorepo Support
Bumpit has built-in support for monorepo versioning. See [examples/config-examples/monorepo.yaml](examples/config-examples/monorepo.yaml) and [examples/workflows/monorepo.yml](examples/workflows/monorepo.yml) for examples.

### Pre-releases
Create beta, alpha, or RC versions with pre-release identifiers. See [examples/workflows/pre-release.yml](examples/workflows/pre-release.yml) for an example.

### Environment Variables
Bumpit supports environment variables in configuration values:
- `${GITHUB_RUN_NUMBER}` - Use in pre-release for build numbers
- `${GITHUB_SHA}` - Use in build metadata for commit hashes
- `${PKG_NAME}` - Use in monorepo setups for package names

### Custom Version Formats
Create your own version format to match your needs:
```yaml
# Standard: v1.2.3
version_format: "{major}.{minor}.{patch}"

# Calendar versioning: v2024.01.1
version_format: "{now:%Y}.{now:%m}.{patch}"

# Custom prefix: release-1.2.3
version_prefix: "release-"
version_format: "{major}.{minor}.{patch}"
```

### Integration Examples
- **Node.js**: `bumpit "npm version ${version}"`
- **Python**: `bumpit "sed -i '' 's/version = .*/version = \"${version}\"/' setup.py"`
- **Gradle**: `bumpit "./gradlew setVersion -Pversion=${version}"`
- **Maven**: `bumpit "mvn versions:set -DnewVersion=${version}"`
- **Cargo**: `bumpit "cargo set-version ${version}"`

## GitHub Action Reference

### Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `version_prefix` | Prefix for version tags (e.g., "v") | No | "v" |
| `version_format` | Format using {major}, {minor}, {patch} | No | "{major}.{minor}.{patch}" |
| `pre_release` | Pre-release identifier (e.g., "alpha.1") | No | "" |
| `build_metadata` | Build metadata to append | No | "" |
| `config_file` | Path to custom config file | No | "" |
| `auto_push` | Whether to push tags automatically | No | "true" |
| `create_release` | Whether to create GitHub release | No | "true" |
| `github_token` | GitHub token for releases | No | ${{ github.token }} |

### Outputs

| Output | Description |
|--------|-------------|
| `version` | The new version number |
| `tag` | The new version tag |
| `previous_version` | The previous version number |
| `is_initial_version` | Whether this is the first version |

### Advanced Usage

```yaml
name: Release
on:
  push:
    branches: [main]

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
        with:
          fetch-depth: 0

      - name: Update version
        id: version
        uses: crazywolf132/bumpit@v1
        with:
          version_prefix: ''
          version_format: '0.{major}.{minor}'
          pre_release: 'beta.1'
          build_metadata: ${{ github.sha }}
          auto_push: true
          create_release: true
          github_token: ${{ secrets.GITHUB_TOKEN }}

      - name: Use version information
        run: |
          echo "New version: ${{ steps.version.outputs.version }}"
          echo "New tag: ${{ steps.version.outputs.tag }}"
          echo "Previous: ${{ steps.version.outputs.previous_version }}"
```

## How It Works

1. **Version Detection**: Finds the latest version tag in your git repository
2. **Commit Analysis**: Reads all commits since that tag
3. **Semantic Analysis**: Analyzes commits using conventional commit format:
   - `BREAKING CHANGE:` or `major:` → Major version bump
   - `feat:` → Minor version bump
   - `fix:`, `chore:`, `docs:`, etc. → Patch version bump
4. **Version Calculation**: Determines the next version using semantic versioning rules
5. **Format Application**: Applies your custom version format
6. **Command Execution**: Runs your specified command with the new version

If no previous version tag exists, it starts from `v0.0.0`.

## License

MIT
