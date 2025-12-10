# Architecture

## Overview

easygh is built using a layered architecture with clear separation of concerns. The application follows Go best practices and design principles.

## Architecture Diagram

```
┌─────────────────────────────────┐
│         CLI Layer               │
│  - Flag parsing                 │
│  - User interaction             │
│  - Command execution            │
└────────────┬────────────────────┘
             │
┌────────────▼────────────────────┐
│      Application Layer          │
│  - Config loading               │
│  - Business logic               │
│  - Error handling               │
│  - Concurrent processing        │
└────────────┬────────────────────┘
             │
┌────────────▼────────────────────┐
│      GitHub API Layer           │
│  - API client                   │
│  - Secrets operations           │
│  - Encryption                   │
└────────────┬────────────────────┘
             │
┌────────────▼────────────────────┐
│      Infrastructure Layer       │
│  - Logging                      │
└─────────────────────────────────┘
```

## Package Structure

### cmd/easygh/

Entry point of the application. Contains the main function and CLI command setup.

### internal/config/

Configuration management:
- `config.go`: Configuration structs and validation
- `loader.go`: YAML file loading and environment variable support

### internal/github/

GitHub API client:
- `client.go`: GitHub client interface and implementation
- `secrets.go`: Secret encryption and API operations

### internal/logger/

Structured logging wrapper around charmbracelet/log.

### internal/cli/

CLI flag parsing and utilities:
- `flags.go`: Flag definitions and parsing helpers


## Design Principles

### Single Responsibility Principle

Each package has a single, well-defined responsibility:
- `config`: Configuration management only
- `github`: GitHub API operations only
- `logger`: Logging only
- `cli`: CLI flag parsing only

### Dependency Injection

Dependencies are injected through constructors, making the code testable:
- `github.NewClient(token)` - creates a GitHub client
- `logger.New(verbose)` - creates a logger
- `config.LoadConfig(path)` - loads configuration

### Interface-Based Design

The GitHub client uses an interface (`github.Client`) which allows for:
- Easy testing with mocks
- Future extensibility (e.g., different implementations)
- Clear contract definition

### Error Handling

Errors are handled explicitly:
- All errors are wrapped with context using `fmt.Errorf` with `%w`
- Errors are propagated up the call stack
- Final error collection and reporting

### Concurrency

Repositories are processed concurrently using goroutines:
- Each repository is processed in its own goroutine
- WaitGroup is used to wait for all goroutines
- Mutex is used for thread-safe error collection

## Data Flow

1. **CLI Parsing**: User provides config file path and flags
2. **Config Loading**: YAML file is loaded and merged with environment variables and CLI flags
3. **Validation**: Configuration is validated
4. **Client Creation**: GitHub client is created with the token
5. **Concurrent Processing**: For each repository:
   - Get public key
   - Encrypt secrets
   - Set secrets via API
6. **Error Collection**: Errors are collected and reported
7. **Result Reporting**: Final status is reported to the user

## Security Considerations

### Secret Encryption

Secrets are encrypted using GitHub's public key encryption:
1. Repository's public key is retrieved
2. Secret is encrypted using NaCl box encryption
3. Encrypted secret is sent to GitHub API

### Token Handling

- Tokens are never logged
- Tokens can be provided via environment variable (more secure than config file)
- Token validation happens at API call time

### Secret Masking

In dry-run mode, secret values are masked in logs (only first 2 and last 2 characters shown).

## Extension Points

### Adding New Secret Sources

To add support for other secret sources (e.g., HashiCorp Vault):
1. Create a new package `internal/secrets/`
2. Define an interface for secret retrieval
3. Implement the interface for each source
4. Update the application layer to use the new source

### Adding New Output Formats

To add support for other output formats:
1. Create a new package `internal/output/`
2. Define an interface for output formatting
3. Implement the interface for each format
4. Update the CLI layer to use the new format

### Adding New GitHub Operations

To add support for other GitHub operations:
1. Add methods to the `github.Client` interface
2. Implement the methods in `githubClient`
3. Update the application layer to use the new operations

## Testing Strategy

### Unit Tests

- Each package has its own test file
- Tests use mocks for external dependencies
- Test fixtures are in `test/fixtures/`

### Integration Tests

- Mock GitHub client is used for integration tests
- Tests verify the full flow without hitting real API
- Mock implementation is in `test/mocks/`

## Build and Deployment

### Makefile

The Makefile provides common build targets:
- `make build`: Build the binary
- `make test`: Run tests
- `make lint`: Run linters
- `make clean`: Clean build artifacts

## Future Improvements

- Rate limiting for GitHub API calls
- Retry logic with exponential backoff
- Progress bar for better UX
- Support for secret deletion
- Support for secret rotation
- Support for multiple GitHub instances (GitHub Enterprise)

