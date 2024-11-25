# Contributing to Bumpit

👋 First off, thanks for taking the time to contribute!

## Quick Start

1. Fork the repository
2. Clone your fork
3. Create a new branch (`git checkout -b feature/amazing-feature`)
4. Make your changes
5. Run tests (`make test`)
6. Commit your changes using [Conventional Commits](https://www.conventionalcommits.org/)
7. Push to your branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## Development Prerequisites

- Go 1.21 or later
- Make

## Development Workflow

1. **Setup**:
   ```bash
   make setup
   ```

2. **Build**:
   ```bash
   make build
   ```

3. **Test**:
   ```bash
   make test
   ```

4. **Lint**:
   ```bash
   make lint
   ```

## Commit Messages

We use [Conventional Commits](https://www.conventionalcommits.org/). This allows us to automatically generate changelogs and determine the next version number.

Examples:
- `feat: add support for custom commit types`
- `fix: handle empty git repositories`
- `docs: improve configuration examples`
- `chore: update dependencies`
- `BREAKING CHANGE: remove deprecated config options`

## Pull Request Process

1. Update the README.md with details of changes if needed
2. Add tests for any new features
3. Update documentation if needed
4. Ensure the test suite passes
5. Update the example workflows if relevant

## Code Style

We use `golangci-lint` for code quality. Run `make lint` before submitting PRs.

Key guidelines:
- Write clear, self-documenting code
- Add comments for complex logic
- Keep functions focused and small
- Write meaningful test cases

## Creating Issues

When creating issues, please:
1. Use a clear and descriptive title
2. Provide detailed reproduction steps
3. Include relevant logs/screenshots
4. List your environment details

## Getting Help

- 💬 Join our [Discord](https://discord.gg/bumpit)
- 📫 Email us at support@bumpit.dev
- 🐛 Report issues on GitHub

## Code of Conduct

This project follows the [Contributor Covenant](https://www.contributor-covenant.org/). Please read [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) before contributing.
