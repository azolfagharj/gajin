# easygh

A CLI tool to manage GitHub Actions secrets across multiple repositories using a YAML configuration file.

## Features

- Manage secrets across multiple repositories with a single configuration file
- Support for YAML configuration with environment variable overrides
- Dry-run mode to preview changes before applying them
- Concurrent processing for faster execution
- Continue on error option for batch operations
- Verbose logging for debugging

## Installation

### From Source

```bash
git clone https://github.com/yourusername/easy_gh_secret.git
cd easy_gh_secret
make build
```

### Binary Release

Download the latest release from the [Releases](https://github.com/yourusername/easy_gh_secret/releases) page.

## Quick Start

1. Create a configuration file (`config.yaml`):

```yaml
github:
  token: GH_TOKEN_WITH_ACTIONS_WRITE  # Or set via GH_TOKEN_WITH_ACTIONS_WRITE env var
  owner: my-org
  repos:
    - repo1
    - repo2

secrets:
  MY_SECRET: "secret-value"
  ANOTHER_SECRET: "another-value"
```

2. Set your GitHub token (if not in config file):

```bash
export GH_TOKEN_WITH_ACTIONS_WRITE=your_github_token
```

3. Run the tool:

```bash
./easygh --config config.yaml
```

## Usage

### Basic Usage

```bash
easygh --config config.yaml
```

### Dry Run Mode

Preview what would be changed without making actual changes:

```bash
easygh --config config.yaml --dry-run
```

### Override Configuration

Override any configuration value from the command line:

```bash
easygh --config config.yaml --token my-token --owner my-org --repo repo1,repo2
```

### Continue on Error

Continue processing other repositories even if one fails:

```bash
easygh --config config.yaml --continue-on-error
```

### Verbose Logging

Enable verbose logging for debugging:

```bash
easygh --config config.yaml --verbose
```

### Show Version

```bash
easygh --version
```

## Configuration

See [examples/config.yaml](examples/config.yaml) for a complete example.

### Configuration File Structure

```yaml
github:
  token: <github-token>  # Optional if GH_TOKEN_WITH_ACTIONS_WRITE is set
  owner: <org-or-username>
  repos:
    - <repo-name-1>
    - <repo-name-2>

secrets:
  <SECRET_NAME>: "<secret-value>"
  <ANOTHER_SECRET>: "<another-value>"
```

### Environment Variables

- `GH_TOKEN_WITH_ACTIONS_WRITE`: GitHub token with `repo` and `actions:write` permissions

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

See the [examples/](examples/) directory for more examples.

## Documentation

- [Usage Guide](docs/USAGE.md) - Detailed usage instructions
- [Architecture](docs/ARCHITECTURE.md) - Architecture and design decisions

## Requirements

- Go 1.22 or later
- GitHub token with `repo` and `actions:write` permissions

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
