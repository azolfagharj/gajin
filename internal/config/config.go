package config

import "fmt"

// Config represents the application configuration.
type Config struct {
	GitHub  GitHubConfig       `yaml:"github"`
	Secrets map[string]string  `yaml:"secrets"`
}

// GitHubConfig contains GitHub-specific configuration.
type GitHubConfig struct {
	Token string   `yaml:"token"`
	Owner string   `yaml:"owner"`
	Repos []string `yaml:"repos"`
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	if c.GitHub.Owner == "" {
		return fmt.Errorf("github.owner is required")
	}

	if len(c.GitHub.Repos) == 0 {
		return fmt.Errorf("at least one repository must be specified in github.repos")
	}

	if c.GitHub.Token == "" {
		return fmt.Errorf("github.token is required (can be set via GH_TOKEN_WITH_ACTIONS_WRITE environment variable)")
	}

	if len(c.Secrets) == 0 {
		return fmt.Errorf("at least one secret must be specified")
	}

	for _, repo := range c.GitHub.Repos {
		if repo == "" {
			return fmt.Errorf("repository name cannot be empty")
		}
	}

	for key, value := range c.Secrets {
		if key == "" {
			return fmt.Errorf("secret key cannot be empty")
		}
		if value == "" {
			return fmt.Errorf("secret value for '%s' cannot be empty", key)
		}
	}

	return nil
}

// ApplyOverrides applies CLI flag overrides to the configuration.
func (c *Config) ApplyOverrides(token, owner string, repos []string) {
	if token != "" {
		c.GitHub.Token = token
	}

	if owner != "" {
		c.GitHub.Owner = owner
	}

	if len(repos) > 0 {
		c.GitHub.Repos = repos
	}
}

