package config

import "fmt"

// Config represents the application configuration.
type Config struct {
	GitHub               GitHubConfig                      `yaml:"github"`
	RepositorySecrets    map[string]string                 `yaml:"repository_secrets"`
	EnvironmentSecrets   map[string]map[string]string      `yaml:"environment_secrets"`
	RepositoryVariables  map[string]string                 `yaml:"repository_variables"`
	EnvironmentVariables map[string]map[string]string      `yaml:"environment_variables"`
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

	// Check if at least one section is specified
	hasRepositorySecrets := len(c.RepositorySecrets) > 0
	hasEnvironmentSecrets := len(c.EnvironmentSecrets) > 0
	hasRepositoryVariables := len(c.RepositoryVariables) > 0
	hasEnvironmentVariables := len(c.EnvironmentVariables) > 0

	if !hasRepositorySecrets && !hasEnvironmentSecrets && !hasRepositoryVariables && !hasEnvironmentVariables {
		return fmt.Errorf("at least one of repository_secrets, environment_secrets, repository_variables, or environment_variables must be specified")
	}

	for _, repo := range c.GitHub.Repos {
		if repo == "" {
			return fmt.Errorf("repository name cannot be empty")
		}
	}

	// Validate repository secrets
	for key, value := range c.RepositorySecrets {
		if key == "" {
			return fmt.Errorf("repository secret key cannot be empty")
		}
		if value == "" {
			return fmt.Errorf("repository secret value for '%s' cannot be empty", key)
		}
	}

	// Validate environment secrets
	for envName, secrets := range c.EnvironmentSecrets {
		if envName == "" {
			return fmt.Errorf("environment name cannot be empty")
		}
		for key, value := range secrets {
			if key == "" {
				return fmt.Errorf("environment secret key cannot be empty for environment '%s'", envName)
			}
			if value == "" {
				return fmt.Errorf("environment secret value for '%s' in environment '%s' cannot be empty", key, envName)
			}
		}
	}

	// Validate repository variables
	for key, value := range c.RepositoryVariables {
		if key == "" {
			return fmt.Errorf("repository variable key cannot be empty")
		}
		if value == "" {
			return fmt.Errorf("repository variable value for '%s' cannot be empty", key)
		}
	}

	// Validate environment variables
	for envName, variables := range c.EnvironmentVariables {
		if envName == "" {
			return fmt.Errorf("environment name cannot be empty")
		}
		for key, value := range variables {
			if key == "" {
				return fmt.Errorf("environment variable key cannot be empty for environment '%s'", envName)
			}
			if value == "" {
				return fmt.Errorf("environment variable value for '%s' in environment '%s' cannot be empty", key, envName)
			}
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

