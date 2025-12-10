# GajIn

**GitHub Actions Secrets & Variables Injector**

A CLI tool to manage GitHub Actions secrets and variables across multiple repositories using a YAML configuration file.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
  - [Binary Release (Recommended)](#binary-release-recommended)
  - [Build from Source](#build-from-source)
- [Usage](#usage)
  - [Basic Usage](#basic-usage)
  - [Dry Run Mode](#dry-run-mode)
  - [Override Configuration](#override-configuration)
  - [Continue on Error](#continue-on-error)
  - [Verbose Logging](#verbose-logging)
  - [Show Version](#show-version)
- [Configuration](#configuration)
  - [Configuration File Structure](#configuration-file-structure)
  - [Environment Variables](#environment-variables)
  - [GitHub Token Permissions](#github-token-permissions)
  - [CLI Flags](#cli-flags)
- [Examples](#examples)
- [Documentation](#documentation)
- [Requirements](#requirements)
- [License](#license)
- [Contributing](#contributing)

## Features

- Manage repository-level and environment-level secrets and variables
- Support for YAML configuration with environment variable overrides
- Dry-run mode to preview changes before applying them
- Concurrent processing for faster execution
- Continue on error option for batch operations
- Verbose logging for debugging
- Automatic encryption for secrets (using GitHub's public key encryption)
- Plaintext storage for variables (no encryption needed)

## Installation

### Binary Release (Recommended)

Download the latest release binary for your system from the [Latest Release](https://github.com/azolfagharj/gajin/releases/latest) page.

**Available binaries:**
- `gajin-darwin-amd64` - macOS (Intel)
- `gajin-darwin-arm64` - macOS (Apple Silicon)
- `gajin-linux-amd64` - Linux (64-bit)
- `gajin-linux-arm64` - Linux (ARM64)
- `gajin-windows-amd64.exe` - Windows (64-bit)

**Quick setup:**

1. Download the binary for your system:
   ```bash
   # Linux (amd64)
   wget https://github.com/azolfagharj/gajin/releases/latest/download/gajin-linux-amd64
   chmod +x gajin-linux-amd64
   mv gajin-linux-amd64 gajin

   # Linux (ARM64)
   wget https://github.com/azolfagharj/gajin/releases/latest/download/gajin-linux-arm64
   chmod +x gajin-linux-arm64
   mv gajin-linux-arm64 gajin

   # macOS (Apple Silicon)
   wget https://github.com/azolfagharj/gajin/releases/latest/download/gajin-darwin-arm64
   chmod +x gajin-darwin-arm64
   mv gajin-darwin-arm64 gajin

   # macOS (Intel)
   wget https://github.com/azolfagharj/gajin/releases/latest/download/gajin-darwin-amd64
   chmod +x gajin-darwin-amd64
   mv gajin-darwin-amd64 gajin

   # Windows (amd64)
   # Download: https://github.com/azolfagharj/gajin/releases/latest/download/gajin-windows-amd64.exe
   # Rename it to gajin.exe
   ```

2. Download or create a configuration file:

   **Option A: Download compact version** (recommended for quick start):
   ```bash
   wget https://github.com/azolfagharj/gajin/releases/latest/download/config.compact.yaml
   mv config.compact.yaml config.yaml
   ```

   **Option B: Create your own** (`config.yaml`):
   ```yaml
   github:
     token: your-github-token-here  # Or use environment variable (see step 3)
     owner: my-org
     repos:
       - repo1
       - repo2

   repository_secrets:
     MY_SECRET: "secret-value"
     ANOTHER_SECRET: "another-value"

   repository_variables:
     LOG_LEVEL: "info"
     DEFAULT_REGION: "us-east-1"

   environment_secrets:
     production:
       PROD_API_KEY: "prod-key"
     staging:
       STAGING_API_KEY: "staging-key"

   environment_variables:
     production:
       DEPLOYMENT_REGION: "us-east-1"
   ```

   **For complete example with detailed comments**, see [examples/config.yaml](https://github.com/azolfagharj/gajin/blob/main/examples/config.yaml) in the repository.

3. Set your GitHub token:

   **Option A: Set in config file** (as shown above)

   **Option B: Use environment variable** (more secure, recommended):
   ```bash
   export GH_TOKEN_WITH_ACTIONS_WRITE=your_github_token
   ```
   If using environment variable, you can remove the `token:` line from config.yaml or leave it as `token: GH_TOKEN_WITH_ACTIONS_WRITE` (it will be read from environment).

   **Note:** If both config file and environment variable are set, the environment variable takes precedence.

4. Run the tool:

```bash
# Linux/macOS
./gajin --config config.yaml

# Windows
gajin.exe --config config.yaml
```

### Build from Source

If you prefer to build from source or need a custom build:

**Prerequisites:**
- Go 1.22 or later

**Build steps:**

```bash
git clone https://github.com/azolfagharj/gajin.git
cd gajin
make build
```

The binary will be created in `bin/gajin`.

**Alternative build command:**

```bash
go build -o bin/gajin ./cmd/gajin
```

## Usage

### Basic Usage

```bash
gajin --config config.yaml
```

### Dry Run Mode

Preview what would be changed without making actual changes:

```bash
gajin --config config.yaml --dry-run
```

### Override Configuration

Override any configuration value from the command line:

```bash
gajin --config config.yaml --token my-token --owner my-org --repo repo1,repo2
```

### Continue on Error

Continue processing other repositories even if one fails:

```bash
gajin --config config.yaml --continue-on-error
```

### Verbose Logging

Enable verbose logging for debugging:

```bash
gajin --config config.yaml --verbose
```

### Show Version

```bash
gajin --version
```

## Configuration

**Configuration Examples:**

- **Compact version** (recommended for quick start): Download `config.compact.yaml` from the [Latest Release](https://github.com/azolfagharj/gajin/releases/latest) page.
- **Complete version** (with detailed comments): See [examples/config.yaml](https://github.com/azolfagharj/gajin/blob/main/examples/config.yaml) in the repository.

### Configuration File Structure

```yaml
github:
  token: <github-token>  # Optional if GH_TOKEN_WITH_ACTIONS_WRITE is set
  owner: <org-or-username>
  repos:
    - <repo-name-1>
    - <repo-name-2>

# Repository-level secrets (encrypted, available to all workflows)
repository_secrets:
  <SECRET_NAME>: "<secret-value>"

# Environment-level secrets (encrypted, available only to workflows using the environment)
environment_secrets:
  <environment-name>:
    <SECRET_NAME>: "<secret-value>"

# Repository-level variables (plaintext, available to all workflows)
repository_variables:
  <VARIABLE_NAME>: "<variable-value>"

# Environment-level variables (plaintext, available only to workflows using the environment)
environment_variables:
  <environment-name>:
    <VARIABLE_NAME>: "<variable-value>"
```

**Note:** At least one of the four sections (`repository_secrets`, `environment_secrets`, `repository_variables`, or `environment_variables`) must be specified.

**Important:** Environments must exist in the repository before setting environment secrets or variables. Create them in Repository Settings > Environments.

### Environment Variables

- `GH_TOKEN_WITH_ACTIONS_WRITE`: GitHub token with appropriate permissions (see below)

### GitHub Token Permissions

**For Fine-grained Personal Access Tokens:**
- **Repository Secrets**: Actions > Secrets: Read and write
- **Repository Variables**: Actions > Variables: Read and write
- **Environment Secrets**: Environments: Read and write (NOT under Actions, it's a separate permission)
- **Environment Variables**: Environments: Read and write (NOT under Actions, it's a separate permission)
- **Metadata**: Read-only (required for all operations)

**For Classic Personal Access Tokens:**
- **All operations**: `repo` scope (provides full access to secrets and variables)
- For public repositories only: `public_repo` scope
- For organization-level operations: `admin:org` scope

For detailed permission setup instructions:
- **Compact version**: Download `config.compact.yaml` from the [Latest Release](https://github.com/azolfagharj/gajin/releases/latest) page.
- **Complete version**: See [examples/config.yaml](https://github.com/azolfagharj/gajin/blob/main/examples/config.yaml) in the repository.

### CLI Flags

- `--config, -c`: Path to configuration file (default: `config.yaml`)
- `--token`: GitHub token (overrides config file)
- `--owner`: GitHub owner/organization (overrides config file)
- `--repo`: Comma-separated list of repositories (overrides config file)
- `--dry-run`: Show what would be done without making changes
- `--continue-on-error`: Continue processing other repositories on error
- `--verbose, -v`: Enable verbose logging
- `--version`: Show version information

## Examples

**Configuration Examples:**

- **Compact version** (recommended for quick start): Download `config.compact.yaml` from the [Latest Release](https://github.com/azolfagharj/gajin/releases/latest) page.
- **Complete version** (with detailed comments): See [examples/config.yaml](https://github.com/azolfagharj/gajin/blob/main/examples/config.yaml) in the repository.

See the [examples/](examples/) directory for more examples.

## Documentation

- [Usage Guide](docs/USAGE.md) - Detailed usage instructions
- [Architecture](docs/ARCHITECTURE.md) - Architecture and design decisions
- [Migration Guide](docs/MIGRATION.md) - How to migrate from the old configuration format

## Requirements

- Go 1.22 or later
- GitHub token with `repo` and `actions:write` permissions

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
