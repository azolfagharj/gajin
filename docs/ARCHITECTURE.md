# Architecture

## Overview

easygh is built using a layered architecture with clear separation of concerns. The application follows Go best practices and design principles.

## Table of Contents

- [Overview](#overview)
- [Architecture Diagram](#architecture-diagram)
- [Package Structure](#package-structure)
  - [cmd/easygh/](#cmdeasygh)
  - [internal/config/](#internalconfig)
  - [internal/github/](#internalgithub)
  - [internal/logger/](#internallogger)
  - [internal/cli/](#internalcli)
- [Design Principles](#design-principles)
  - [Single Responsibility Principle](#single-responsibility-principle)
  - [Dependency Injection](#dependency-injection)
  - [Interface-Based Design](#interface-based-design)
  - [Error Handling](#error-handling)
  - [Concurrency](#concurrency)
- [Data Flow](#data-flow)
- [Security Considerations](#security-considerations)
  - [Secret Encryption](#secret-encryption)
  - [Token Handling](#token-handling)
  - [Secret Masking](#secret-masking)
- [Extension Points](#extension-points)
  - [Adding New Secret Sources](#adding-new-secret-sources)
  - [Adding New Output Formats](#adding-new-output-formats)
  - [Adding New GitHub Operations](#adding-new-github-operations)
- [Testing Strategy](#testing-strategy)
  - [Unit Tests](#unit-tests)
  - [Integration Tests](#integration-tests)
- [Build and Deployment](#build-and-deployment)
  - [Binary Releases](#binary-releases)
  - [Building from Source](#building-from-source)
  - [Makefile](#makefile)
  - [Cross-platform Builds](#cross-platform-builds)
- [Supported Operations](#supported-operations)
  - [Repository Secrets](#repository-secrets)
  - [Environment Secrets](#environment-secrets)
  - [Repository Variables](#repository-variables)
  - [Environment Variables](#environment-variables)
- [Future Improvements](#future-improvements)

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
- `secrets.go`: Secret encryption and API operations for repository and environment secrets
- `errors.go`: Custom error types for better error handling

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
   - Get repository ID (for environment operations)
   - Process repository secrets (get public key, encrypt, set)
   - Process environment secrets (get environment public key, encrypt, set)
   - Process repository variables (set directly, no encryption)
   - Process environment variables (set directly, no encryption)
6. **Error Collection**: Errors are collected and reported with context
7. **Result Reporting**: Final status is reported to the user

## Security Considerations

### Secret Encryption

Secrets are encrypted using GitHub's public key encryption following the LibSodium sealed box format (`crypto_box_seal`):

1. Repository's public key is retrieved from GitHub API
2. An ephemeral key pair is generated for each encryption
3. A nonce is derived using **BLAKE2b with 24-byte output** (NOT truncated from 64 bytes):
   ```
   nonce = BLAKE2b(ephemeral_pk || recipient_pk, digest_size=24)
   ```
4. Secret is encrypted using NaCl box encryption with the derived nonce
5. Encrypted data format: `[ephemeral public key (32 bytes)][encrypted ciphertext + MAC (16 bytes)]`
6. The encrypted value is base64 encoded and sent to GitHub API along with the key ID

**Important:** The nonce derivation uses `BLAKE2b` with `digest_size=24`, which is different from computing `BLAKE2b-512` and taking the first 24 bytes. This distinction is critical for compatibility with LibSodium's `crypto_box_seal`.

This format ensures that each encryption produces different ciphertext even for the same plaintext, providing forward secrecy.

### Token Handling

- Tokens are never logged
- Tokens can be provided via environment variable (more secure than config file)
- Token validation happens at API call time

### Secret Masking

In dry-run mode and logs, secret values are masked (only first 2 and last 2 characters shown). Variable values are shown in full since they are not sensitive.

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

### Binary Releases

Pre-built binaries are available for all supported platforms in each [release](https://github.com/azolfagharj/easy_gh_secret/releases):
- `easygh-darwin-amd64` - macOS (Intel)
- `easygh-darwin-arm64` - macOS (Apple Silicon)
- `easygh-linux-amd64` - Linux (64-bit)
- `easygh-linux-arm64` - Linux (ARM64)
- `easygh-windows-amd64.exe` - Windows (64-bit)

**Recommended:** Use the binary release for the fastest setup. Download from [Latest Release](https://github.com/azolfagharj/easy_gh_secret/releases/latest).

### Building from Source

### Makefile

The Makefile provides common build targets:
- `make build`: Build the binary for current platform
- `make test`: Run tests
- `make lint`: Run linters
- `make clean`: Clean build artifacts

### Cross-platform Builds

To build for different platforms:

```bash
# Linux (amd64)
GOOS=linux GOARCH=amd64 go build -o bin/easygh-linux-amd64 ./cmd/easygh

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o bin/easygh-darwin-arm64 ./cmd/easygh

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o bin/easygh-darwin-amd64 ./cmd/easygh

# Linux (ARM64)
GOOS=linux GOARCH=arm64 go build -o bin/easygh-linux-arm64 ./cmd/easygh

# Windows
GOOS=windows GOARCH=amd64 go build -o bin/easygh-windows-amd64.exe ./cmd/easygh
```

## Supported Operations

### Repository Secrets
- Create/update repository secrets with automatic encryption
- Read secret metadata (values cannot be retrieved)

### Environment Secrets
- Create/update environment secrets with automatic encryption
- Read secret metadata (values cannot be retrieved)
- Requires environment to exist in repository

### Repository Variables
- Create/update repository variables (plaintext)
- Read variable values (including the actual value)

### Environment Variables
- Create/update environment variables (plaintext)
- Read variable values (including the actual value)
- Requires environment to exist in repository

## Future Improvements

- Rate limiting for GitHub API calls
- Retry logic with exponential backoff
- Progress bar for better UX
- Support for secret/variable deletion
- Support for secret rotation
- Support for multiple GitHub instances (GitHub Enterprise)
- Bulk operations for better performance

