# Usage Guide

## Table of Contents

- [Installation](#installation)
- [Configuration](#configuration)
- [Basic Usage](#basic-usage)
- [Advanced Usage](#advanced-usage)
- [Troubleshooting](#troubleshooting)

## Installation

### Binary Release (Recommended)

Download the latest release binary for your system from the [Latest Release](https://github.com/azolfagharj/easy_gh_secret/releases/latest) page.

**Available binaries:**
- `easygh-darwin-amd64` - macOS (Intel)
- `easygh-darwin-arm64` - macOS (Apple Silicon)
- `easygh-linux-amd64` - Linux (64-bit)
- `easygh-linux-arm64` - Linux (ARM64)
- `easygh-windows-amd64.exe` - Windows (64-bit)

**Quick setup:**

1. Download the binary for your system:
   ```bash
   # Linux (amd64)
   wget https://github.com/azolfagharj/easy_gh_secret/releases/latest/download/easygh-linux-amd64
   chmod +x easygh-linux-amd64
   mv easygh-linux-amd64 easygh

   # Linux (ARM64)
   wget https://github.com/azolfagharj/easy_gh_secret/releases/latest/download/easygh-linux-arm64
   chmod +x easygh-linux-arm64
   mv easygh-linux-arm64 easygh

   # macOS (Apple Silicon)
   wget https://github.com/azolfagharj/easy_gh_secret/releases/latest/download/easygh-darwin-arm64
   chmod +x easygh-darwin-arm64
   mv easygh-darwin-arm64 easygh

   # macOS (Intel)
   wget https://github.com/azolfagharj/easy_gh_secret/releases/latest/download/easygh-darwin-amd64
   chmod +x easygh-darwin-amd64
   mv easygh-darwin-amd64 easygh

   # Windows (amd64)
   # Download: https://github.com/azolfagharj/easy_gh_secret/releases/latest/download/easygh-windows-amd64.exe
   # Rename it to easygh.exe
   ```

2. Verify installation:
   ```bash
   ./easygh --version
   ```

### Build from Source

If you prefer to build from source or need a custom build:

**Prerequisites:**
- Go 1.22 or later
- GitHub token with appropriate permissions (see [GitHub Token Permissions](#github-token-permissions))

**Build steps:**

```bash
git clone https://github.com/azolfagharj/easy_gh_secret.git
cd easy_gh_secret
make build
```

The binary will be created in `bin/easygh`.

**Alternative build command:**

```bash
go build -o bin/easygh ./cmd/easygh
```

**Cross-platform builds:**

```bash
# Build for Linux (amd64)
GOOS=linux GOARCH=amd64 go build -o bin/easygh-linux-amd64 ./cmd/easygh

# Build for macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o bin/easygh-darwin-arm64 ./cmd/easygh

# Build for macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o bin/easygh-darwin-amd64 ./cmd/easygh

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o bin/easygh-windows-amd64.exe ./cmd/easygh
```

## Configuration

### Configuration File

Create a YAML configuration file (default: `config.yaml`):

```yaml
github:
  token: your-github-token  # Optional if GH_TOKEN_WITH_ACTIONS_WRITE is set
  owner: my-organization
  repos:
    - repository-1
    - repository-2
    - repository-3

# Repository-level secrets (encrypted, available to all workflows)
repository_secrets:
  DATABASE_URL: "postgresql://localhost:5432/mydb"
  API_KEY: "sk-1234567890abcdef"
  WEBHOOK_SECRET: "whsec_1234567890"

# Environment-level secrets (encrypted, available only to workflows using the environment)
environment_secrets:
  production:
    PROD_DATABASE_URL: "postgresql://prod.db.example.com:5432/mydb"
  staging:
    STAGING_DATABASE_URL: "postgresql://staging.db.example.com:5432/mydb"

# Repository-level variables (plaintext, available to all workflows)
repository_variables:
  LOG_LEVEL: "info"
  DEFAULT_REGION: "us-east-1"

# Environment-level variables (plaintext, available only to workflows using the environment)
environment_variables:
  production:
    DEPLOYMENT_REGION: "us-east-1"
  staging:
    DEPLOYMENT_REGION: "us-west-2"
```

**Note:** At least one of the four sections must be specified. Environments must exist in the repository before setting environment secrets or variables.

### Environment Variables

You can set the GitHub token via environment variable instead of in the config file:

```bash
export GH_TOKEN_WITH_ACTIONS_WRITE=your-github-token
```

If both are set, the environment variable takes precedence.

### GitHub Token Permissions

**Fine-grained Personal Access Tokens (Recommended):**

Your GitHub token needs specific permissions based on what you want to manage:

- **Repository Secrets**:
  - Repository permissions > Actions > **Secrets**: Read and write

- **Repository Variables**:
  - Repository permissions > Actions > **Variables**: Read and write

- **Environment Secrets**:
  - Repository permissions > **Environments**: Read and write
  - IMPORTANT: Environment secrets require Environments permission (NOT under Actions)

- **Environment Variables**:
  - Repository permissions > **Environments**: Read and write
  - IMPORTANT: Environment variables require Environments permission (NOT under Actions)

- **Required for all operations**:
  - Repository permissions > **Metadata**: Read-only (required for API access)

**Classic Personal Access Tokens:**

- **All operations**: `repo` scope
  - This single scope provides full access to:
    - Repository secrets and variables
    - Environment secrets and variables
  - Note: Classic tokens provide broader permissions than needed

For detailed setup instructions, see the comments in [examples/config.yaml](../../examples/config.yaml).

## Basic Usage

### Set Secrets and Variables

```bash
easygh --config config.yaml
```

This will set all secrets and variables defined in the configuration file to all specified repositories. The tool processes:
- Repository secrets and variables (set for all repositories)
- Environment secrets and variables (set for each environment in each repository)

### Dry Run

Preview changes before applying them:

```bash
easygh --config config.yaml --dry-run
```

The dry-run mode will show:
- Which secrets would be created
- Which secrets would be updated (with existing metadata)
- No actual changes will be made

## Advanced Usage

### Override Configuration Values

Override any configuration value from the command line:

```bash
# Override token
easygh --config config.yaml --token different-token

# Override owner
easygh --config config.yaml --owner different-org

# Override repositories
easygh --config config.yaml --repo repo1,repo2,repo3

# Combine multiple overrides
easygh --config config.yaml --owner my-org --repo repo1,repo2 --token my-token
```

### Continue on Error

By default, the tool stops on the first error. To continue processing other repositories:

```bash
easygh --config config.yaml --continue-on-error
```

All errors will be collected and displayed at the end.

### Verbose Logging

Enable verbose logging for debugging:

```bash
easygh --config config.yaml --verbose
```

This will show detailed information about each operation.

### Custom Config File Path

```bash
easygh --config /path/to/my-config.yaml
```

Or use the short form:

```bash
easygh -c /path/to/my-config.yaml
```

## Troubleshooting

### Common Issues

#### "github.token is required"

**Solution**: Set the token either in the config file or via `GH_TOKEN_WITH_ACTIONS_WRITE` environment variable.

```bash
export GH_TOKEN_WITH_ACTIONS_WRITE=your-token
```

#### "failed to get public key"

**Possible causes**:
- Invalid token
- Token doesn't have `repo` permission
- Repository doesn't exist or you don't have access

**Solution**: Verify your token has the correct permissions and the repository exists.

#### "failed to create or update secret"

**Possible causes**:
- Token doesn't have `actions:write` permission
- Repository doesn't have GitHub Actions enabled

**Solution**: Ensure your token has `actions:write` permission and Actions is enabled for the repository.

#### "validation failed due to an improperly encrypted secret"

**Possible causes**:
- Encryption format mismatch (using incorrect nonce derivation)
- Invalid or corrupted public key

**Solution**: This is typically a bug in the encryption implementation. Make sure you're using the latest version of easygh which correctly implements LibSodium sealed box encryption.

### Debugging

Use verbose mode to see detailed error messages:

```bash
easygh --config config.yaml --verbose
```

### Rate Limiting

GitHub API has rate limits. If you encounter rate limiting:
- The tool processes repositories concurrently, which may hit rate limits
- Consider processing repositories in smaller batches
- Use a token with higher rate limits (GitHub App tokens have higher limits)

## Examples

### Example 1: Quick Start with Binary Release

1. **Download the binary** for your system from [Latest Release](https://github.com/azolfagharj/easy_gh_secret/releases/latest):
   ```bash
   # Linux (amd64)
   wget https://github.com/azolfagharj/easy_gh_secret/releases/latest/download/easygh-linux-amd64
   chmod +x easygh-linux-amd64
   mv easygh-linux-amd64 easygh

   # Linux (ARM64)
   wget https://github.com/azolfagharj/easy_gh_secret/releases/latest/download/easygh-linux-arm64
   chmod +x easygh-linux-arm64
   mv easygh-linux-arm64 easygh

   # macOS (Apple Silicon)
   wget https://github.com/azolfagharj/easy_gh_secret/releases/latest/download/easygh-darwin-arm64
   chmod +x easygh-darwin-arm64
   mv easygh-darwin-arm64 easygh

   # macOS (Intel)
   wget https://github.com/azolfagharj/easy_gh_secret/releases/latest/download/easygh-darwin-amd64
   chmod +x easygh-darwin-amd64
   mv easygh-darwin-amd64 easygh

   # Windows: Download easygh-windows-amd64.exe and rename to easygh.exe
   ```

2. **Create configuration file** (`config.yaml`):
   ```yaml
   github:
     token: your-github-token-here  # Or use environment variable (see step 3)
     owner: my-org
     repos:
       - frontend
       - backend
       - api

   repository_secrets:
     DATABASE_URL: "postgresql://prod.db.example.com:5432/mydb"
     REDIS_URL: "redis://prod.redis.example.com:6379"
   ```

3. **Set GitHub token** (if not set in config file):

   **Option A: Already set in config.yaml** (as shown above)

   **Option B: Use environment variable** (more secure, recommended):
   ```bash
   export GH_TOKEN_WITH_ACTIONS_WRITE=ghp_xxxxxxxxxxxx
   ```
   If using environment variable, you can remove the `token:` line from config.yaml or leave it as `token: GH_TOKEN_WITH_ACTIONS_WRITE`. The environment variable will take precedence.

4. **Run the tool**:
   ```bash
   # Linux/macOS
   ./easygh --config config.yaml

   # Windows
   easygh.exe --config config.yaml
   ```

### Example 4: Dry Run Before Applying

```bash
easygh --config config.yaml --dry-run
```

This will show what would be set without making actual changes. Output includes:
- Repository secrets and variables
- Environment secrets and variables (with environment name)

### Example 5: Override Repositories

```bash
easygh --config config.yaml --repo frontend,backend
```

This will only set secrets and variables to `frontend` and `backend` repositories, ignoring the repos in the config file.

## Differences Between Secrets and Variables

### Secrets
- **Encrypted**: Values are encrypted using GitHub's public key encryption before being sent to the API
- **Masked**: Values are masked in logs and workflow outputs (only first 2 and last 2 characters shown)
- **Not readable**: Once set, secret values cannot be retrieved via API (only metadata)
- **Use for**: Passwords, API keys, tokens, database credentials, and other sensitive data

### Variables
- **Plaintext**: Values are stored as plaintext (no encryption)
- **Visible**: Values can be read back via API and shown in logs
- **Use for**: Non-sensitive configuration values like regions, log levels, feature flags, deployment settings

## Environment Requirements

Before setting environment secrets or variables, you must create the environments in GitHub:

1. Go to your repository on GitHub
2. Navigate to **Settings** > **Environments**
3. Click **"New environment"**
4. Enter the environment name (e.g., "production", "staging")
5. Configure protection rules if needed
6. Click **"Configure environment"**

If you try to set environment secrets or variables for a non-existent environment, you'll get an error message indicating that the environment must be created first.

