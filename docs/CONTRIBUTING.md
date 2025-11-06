# Contributing to Onlineliga Discord Bot

Thank you for considering contributing to the Onlineliga Discord Bot! This document provides guidelines and
instructions for contributing to the project.

## Code of Conduct

- Be respectful and inclusive
- Provide constructive feedback
- Focus on what is best for the community
- Show empathy towards other contributors

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports:

1. Check existing issues for duplicates
1. Test with the latest version
1. Collect relevant information (logs, environment, steps to reproduce)

**Bug Report Template:**

````markdown
**Describe the bug**
A clear description of what the bug is.

**To Reproduce**
Steps to reproduce:

1. Command used: `/results location:.de users:12345`
2. Expected behavior
3. Actual behavior

**Environment:**

- Deployment method (Docker/Binary/Compose)
- Go version (if building from source)
- Platform (OS, architecture)

**Logs:**

```text

Paste relevant logs here

```text

**Additional context**
Any other context about the problem.

```text

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues:

1. **Use a clear title** - Describe the enhancement concisely
2. **Provide detailed description** - Explain why this would be useful
3. **Include examples** - Show how the feature would work
4. **Consider alternatives** - Discuss different approaches

### Pull Requests

1. **Fork the repository**
2. **Create a feature branch** from `main`
3. **Make your changes**
4. **Write/update tests**
5. **Update documentation**
6. **Submit pull request**

## Development Setup

See [DEVELOPMENT.md](DEVELOPMENT.md) for detailed setup instructions.

**Quick Start:**

```bash
git clone https://github.com/tripleawwy/onlineliga_discord_bot.git
cd onlineliga_discord_bot
go mod download
go test ./...
go build -o bin/ol_discord_bot .

```text

## Coding Guidelines

### Go Style

Follow standard Go conventions:

**Formatting:**

```bash
go fmt ./...
gofumpt -w .  # Optional: stricter formatting

```text

**Linting:**

```bash
go vet ./...
golangci-lint run

```text

**Naming Conventions:**

- Exported: `PascalCase` (e.g., `ScrapeResults`)
- Unexported: `camelCase` (e.g., `parseJSON`)
- Files: `snake_case.go` (e.g., `scraper.go`, `parse_result.go`)
- Test files: `*_test.go` (e.g., `scraper_test.go`)

### Code Organization

**Package Structure:**

```go
// Package comment describing purpose
package packagename

import (
    // Standard library
    "fmt"
    "log"

    // External packages
    "github.com/sirupsen/logrus"

    // Internal packages
    "github.com/tripleawwy/onlineliga_discord_bot/internal/parse"
)

// Exported type with comment
type PublicType struct {
    Field string
}

// Exported function with comment
func PublicFunction() {
    // Implementation
}

// unexported helper
func helperFunction() {
    // Implementation
}

```text

### Documentation

**Package Documentation:**

```go
// Package scraper provides HTTP data fetching for Onlineliga API.
//
// The scraper package abstracts API communication and handles
// error resilience and multi-user coordination.
package scraper

```text

**Function Documentation:**

```go
// ScrapeMatchResults fetches match results for multiple users.
//
// It processes users sequentially and continues on individual failures,
// returning partial results. Errors are logged but not propagated.
//
// Parameters:
//   - userIDs: Slice of user ID strings
//   - baseURL: Onlineliga base URL (e.g., "https://www.onlineliga.de")
//
// Returns:
//   - Array of successfully fetched MatchResult structs
func (s *Scraper) ScrapeMatchResults(userIDs []string, baseURL string) []parse.MatchResult {
    // Implementation
}

```text

### Error Handling

**Pattern:**

```go
// Return errors, don't panic
result, err := someOperation()
if err != nil {
    return nil, fmt.Errorf("operation failed: %w", err)
}

// Log with context
if err != nil {
    logger.WithError(err).WithFields(logrus.Fields{
        "userID": userID,
        "url":    url,
    }).Error("Failed to fetch data")
}

// Resource cleanup
defer func() {
    if err := resource.Close(); err != nil {
        logger.WithError(err).Warn("Cleanup failed")
    }
}()

```text

### Testing

**Test Organization:**

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        expected OutputType
        wantErr  bool
    }{
        {
            name:     "valid input",
            input:    validInput,
            expected: expectedOutput,
            wantErr:  false,
        },
        {
            name:    "invalid input",
            input:   invalidInput,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := FunctionName(tt.input)

            if (err != nil) != tt.wantErr {
                t.Errorf("FunctionName() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if !reflect.DeepEqual(result, tt.expected) {
                t.Errorf("FunctionName() = %v, want %v", result, tt.expected)
            }
        })
    }
}

```text

**Test Helpers:**

```go
// Test fixtures in internal/*/files/
testData, err := os.ReadFile("files/test_data.json")
if err != nil {
    t.Fatal(err)
}

// HTTP mocks
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"response": "data"}`))
}))
defer server.Close()

```text

## Commit Guidelines

### Commit Message Format

```text

<type>: <subject>

<body>

<footer>

```text

**Types:**

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `test`: Adding/updating tests
- `refactor`: Code refactoring
- `perf`: Performance improvement
- `chore`: Maintenance tasks

**Examples:**

```text

feat: add support for multiple regions

Implement region selection for Onlineliga API endpoints.
Users can now specify .de, .co.uk, .at, or .ch regions.

Closes #42

```text

```text

fix: correct badge rendering position

Badges were offset by 5px due to incorrect y-coordinate
calculation. This commit fixes the positioning logic.

Fixes #38

```text

```text

docs: update API documentation for scraper package

Add examples for concurrent scraping and error handling.
Clarify return types and error conditions.

```text

### Commit Best Practices

- **One logical change per commit**
- **Write clear, descriptive messages**
- **Use present tense** ("add feature" not "added feature")
- **Reference issues** when applicable
- **Keep commits atomic** and reversible

## Pull Request Process

### Before Submitting

1. **Update your fork:**

   ```bash
   git remote add upstream https://github.com/tripleawwy/onlineliga_discord_bot.git
   git fetch upstream
   git rebase upstream/main

   ```text

2. **Run tests:**

   ```bash
   go test ./...

   ```text

3. **Format code:**

   ```bash
   go fmt ./...

   ```text

4. **Lint code:**

   ```bash
   go vet ./...

   ```text

5. **Build successfully:**

   ```bash
   go build -o bin/ol_discord_bot .

   ```text

6. **Update documentation:**
   - Update README if user-facing changes
   - Update API docs if package changes
   - Add/update code comments

### PR Template

```markdown
## Description
Brief description of changes.

## Motivation
Why are these changes needed?

## Changes
- List of specific changes
- Another change
- Yet another change

## Testing
How were these changes tested?

## Screenshots (if applicable)
Add screenshots for UI changes.

## Checklist
- [ ] Tests pass locally
- [ ] Code follows project style
- [ ] Documentation updated
- [ ] No breaking changes (or documented)
- [ ] Commits are atomic and well-described

```text

### Review Process

- Maintainers will review within 1-3 days
- Address feedback promptly
- Maintain civil and constructive discussion
- Be patient and understanding

## Testing Guidelines

### Unit Tests

**Coverage Target:** 70%+

```bash
# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

```text

**What to Test:**

- Public API functions
- Error handling paths
- Edge cases and boundary conditions
- Data transformation logic

**Test Data:**
Store fixtures in `internal/*/files/`:

```text

internal/parse/files/
├── parse_test.json
├── parse_error.json
└── parse_edge_case.json

```text

### Integration Tests

```go
// +build integration

func TestIntegration_ScrapeAPI(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }

    // Test against real API
    scraper := scraper.NewScraper(logger)
    results := scraper.ScrapeMatchResults([]string{"12345"}, "https://www.onlineliga.de")

    if len(results) == 0 {
        t.Error("Expected results from API")
    }
}

```text

Run: `go test -tags=integration ./...`

## Documentation Guidelines

### When to Update Docs

**Always update documentation for:**

- New features or commands
- API changes (function signatures, behavior)
- Configuration changes
- Breaking changes
- New environment variables

**Documentation Locations:**

- User-facing: `README.md`
- Architecture: `ARCHITECTURE.md`
- Development: `DEVELOPMENT.md`
- Deployment: `DEPLOYMENT.md`
- API: `docs/api/<package>.md`

### Documentation Style

**Code Examples:**

```go
// Always include imports
import (
    "github.com/tripleawwy/onlineliga_discord_bot/internal/scraper"
)

// Show complete, runnable examples
func example() {
    scraper := scraper.NewScraper(logger)
    results := scraper.ScrapeMatchResults(userIDs, baseURL)
    // Use results...
}

```text

**Error Handling:**
Always show error handling in examples:

```go
result, err := operation()
if err != nil {
    log.Fatalf("Operation failed: %v", err)
}

```text

## Release Process

### Version Numbering

Follow [Semantic Versioning](https://semver.org/):

- **MAJOR**: Breaking changes
- **MINOR**: New features (backwards compatible)
- **PATCH**: Bug fixes

### Release Checklist

1. Update version in documentation
2. Update CHANGELOG.md
3. Tag release: `git tag -a v1.0.0 -m "Release v1.0.0"`
4. Push tag: `git push origin v1.0.0`
5. CI automatically builds and publishes Docker image
6. Create GitHub release with notes

## Community

### Getting Help

- **Documentation**: Start with README and guides
- **Issues**: Search existing issues
- **Discussions**: Use GitHub Discussions for questions
- **Pull Requests**: Ask questions in PR comments

### Recognition

Contributors are recognized in:

- GitHub Contributors list
- Release notes
- README acknowledgments

## License

By contributing, you agree that your contributions will be licensed under the project's MIT License.

## Questions?

If you have questions about contributing:

1. Check this guide and other documentation
2. Search existing issues and discussions
3. Open a new discussion if needed

---

**Thank you for contributing!** 🎉

Your efforts help make this project better for the Onlineliga community.
````
