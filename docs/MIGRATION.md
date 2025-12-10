# Migration Guide

This guide helps you migrate from the old configuration format to the new format that supports repository and environment-level secrets and variables.

## Table of Contents

- [Breaking Changes](#breaking-changes)
- [Migration Steps](#migration-steps)
  - [Step 1: Update Configuration Structure](#step-1-update-configuration-structure)
  - [Step 2: Add New Sections (Optional)](#step-2-add-new-sections-optional)
  - [Step 3: Verify Token Permissions](#step-3-verify-token-permissions)
  - [Step 4: Create Environments (If Using Environment Secrets/Variables)](#step-4-create-environments-if-using-environment-secretsvariables)
  - [Step 5: Verify Migration](#step-5-verify-migration)
- [Differences Between Secrets and Variables](#differences-between-secrets-and-variables)
- [Examples](#examples)
  - [Example 1: Simple Migration](#example-1-simple-migration)
  - [Example 2: Adding Environment Support](#example-2-adding-environment-support)
  - [Example 3: Using Variables](#example-3-using-variables)
- [Troubleshooting](#troubleshooting)
- [Need Help?](#need-help)

## Breaking Changes

The old `secrets:` key has been removed. You must migrate to the new structure.

## Migration Steps

### Step 1: Update Configuration Structure

**Old Format:**
```yaml
github:
  token: your-token
  owner: my-org
  repos:
    - repo1
    - repo2

secrets:
  SECRET1: "value1"
  SECRET2: "value2"
```

**New Format:**
```yaml
github:
  token: your-token
  owner: my-org
  repos:
    - repo1
    - repo2

repository_secrets:
  SECRET1: "value1"
  SECRET2: "value2"
```

Simply rename `secrets:` to `repository_secrets:`.

### Step 2: Add New Sections (Optional)

You can now add environment secrets and variables:

```yaml
github:
  token: your-token
  owner: my-org
  repos:
    - repo1
    - repo2

repository_secrets:
  SECRET1: "value1"
  SECRET2: "value2"

# New: Environment secrets
environment_secrets:
  production:
    PROD_SECRET: "prod-value"
  staging:
    STAGING_SECRET: "staging-value"

# New: Repository variables
repository_variables:
  LOG_LEVEL: "info"
  DEFAULT_REGION: "us-east-1"

# New: Environment variables
environment_variables:
  production:
    DEPLOYMENT_REGION: "us-east-1"
  staging:
    DEPLOYMENT_REGION: "us-west-2"
```

### Step 3: Verify Token Permissions

Make sure your GitHub token has the required permissions:

**For Fine-grained Tokens:**
- **Repository Secrets**: Actions > Secrets: Read and write
- **Repository Variables**: Actions > Variables: Read and write
- **Environment Secrets**: Environments: Read and write (NOT under Actions, it's a separate permission)
- **Environment Variables**: Environments: Read and write (NOT under Actions, it's a separate permission)
- **Metadata**: Read-only (required for all operations)

**For Classic Tokens:**
- **All operations**: `repo` scope

For detailed permission setup instructions:
- **Compact version**: Download `config.compact.yaml` from the [Latest Release](https://github.com/azolfagharj/gajin/releases/latest) page.
- **Complete version**: See [examples/config.yaml](https://github.com/azolfagharj/gajin/blob/main/examples/config.yaml) in the repository.

### Step 4: Create Environments (If Using Environment Secrets/Variables)

Before setting environment secrets or variables, you must create the environments in GitHub:

1. Go to your repository on GitHub
2. Navigate to Settings > Environments
3. Click "New environment"
4. Enter the environment name (e.g., "production", "staging")
5. Configure protection rules if needed
6. Click "Configure environment"

Repeat for each environment you want to use.

### Step 5: Verify Migration

Run the tool with `--dry-run` to verify your configuration:

```bash
gajin --config config.yaml --dry-run
```

This will show what would be set without making actual changes.

## Differences Between Secrets and Variables

### Secrets
- **Encrypted**: Values are encrypted using GitHub's public key encryption
- **Masked**: Values are masked in logs and workflow outputs
- **Not readable**: Once set, secret values cannot be retrieved via API
- **Use for**: Passwords, API keys, tokens, and other sensitive data

### Variables
- **Plaintext**: Values are stored as plaintext
- **Visible**: Values can be read back via API and shown in logs
- **Use for**: Non-sensitive configuration values like regions, log levels, feature flags

## Examples

### Example 1: Simple Migration

**Before:**
```yaml
secrets:
  DATABASE_URL: "postgresql://..."
  API_KEY: "sk-..."
```

**After:**
```yaml
repository_secrets:
  DATABASE_URL: "postgresql://..."
  API_KEY: "sk-..."
```

### Example 2: Adding Environment Support

**Before:**
```yaml
secrets:
  DATABASE_URL: "postgresql://prod..."
```

**After:**
```yaml
repository_secrets:
  # Shared secrets
  SHARED_SECRET: "value"

environment_secrets:
  production:
    DATABASE_URL: "postgresql://prod..."
  staging:
    DATABASE_URL: "postgresql://staging..."
```

### Example 3: Using Variables

```yaml
repository_variables:
  DEFAULT_REGION: "us-east-1"
  LOG_LEVEL: "info"

environment_variables:
  production:
    DEPLOYMENT_REGION: "us-east-1"
    ENABLE_MONITORING: "true"
  staging:
    DEPLOYMENT_REGION: "us-west-2"
    ENABLE_MONITORING: "false"
```

## Troubleshooting

### Error: "environment 'production' not found"

**Solution**: Create the environment in GitHub repository settings first. See Step 3 above.

### Error: "at least one of repository_secrets, environment_secrets, repository_variables, or environment_variables must be specified"

**Solution**: Make sure you've migrated from `secrets:` to `repository_secrets:` or added one of the new sections.

### Old config still works

If you're using an old version of gajin, update to the latest version first, then migrate your configuration.

## Need Help?

If you encounter issues during migration, please:
1. Check the [Usage Guide](USAGE.md) for detailed examples
2. Review the [Architecture](ARCHITECTURE.md) documentation
3. Open an issue on GitHub with your configuration (redact sensitive values)

