# Usage Guide

## Table of Contents

- [Installation](#installation)
- [Configuration](#configuration)
- [Basic Usage](#basic-usage)
- [Advanced Usage](#advanced-usage)
- [Troubleshooting](#troubleshooting)

## Installation

### Prerequisites

- Go 1.22 or later
- GitHub token with `repo` and `actions:write` permissions

### Build from Source

```bash
git clone https://github.com/yourusername/easy_gh_secret.git
cd easy_gh_secret
make build
```

The binary will be created in `bin/easygh`.

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

secrets:
  DATABASE_URL: "postgresql://localhost:5432/mydb"
  API_KEY: "sk-1234567890abcdef"
  WEBHOOK_SECRET: "whsec_1234567890"
```

### Environment Variables

You can set the GitHub token via environment variable instead of in the config file:

```bash
export GH_TOKEN_WITH_ACTIONS_WRITE=your-github-token
```

If both are set, the environment variable takes precedence.

### GitHub Token Permissions

Your GitHub token needs the following permissions:
- `repo` (Full control of private repositories)
- `actions:write` (Write access to GitHub Actions)

## Basic Usage

### Set Secrets

```bash
easygh --config config.yaml
```

This will set all secrets defined in the configuration file to all specified repositories.

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

### Example 1: Set Secrets to Multiple Repositories

```yaml
# config.yaml
github:
  owner: my-org
  repos:
    - frontend
    - backend
    - api

secrets:
  DATABASE_URL: "postgresql://prod.db.example.com:5432/mydb"
  REDIS_URL: "redis://prod.redis.example.com:6379"
```

```bash
export GH_TOKEN_WITH_ACTIONS_WRITE=ghp_xxxxxxxxxxxx
easygh --config config.yaml
```

### Example 2: Dry Run Before Applying

```bash
easygh --config config.yaml --dry-run
```

### Example 3: Override Repositories

```bash
easygh --config config.yaml --repo frontend,backend
```

This will only set secrets to `frontend` and `backend` repositories, ignoring the repos in the config file.

