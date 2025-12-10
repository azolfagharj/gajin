package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: &Config{
				GitHub: GitHubConfig{
					Token: "test-token",
					Owner: "test-org",
					Repos: []string{"repo1", "repo2"},
				},
				Secrets: map[string]string{
					"SECRET1": "value1",
					"SECRET2": "value2",
				},
			},
			wantErr: false,
		},
		{
			name: "missing owner",
			config: &Config{
				GitHub: GitHubConfig{
					Token: "test-token",
					Owner: "",
					Repos: []string{"repo1"},
				},
				Secrets: map[string]string{"SECRET1": "value1"},
			},
			wantErr: true,
			errMsg:  "github.owner is required",
		},
		{
			name: "missing repos",
			config: &Config{
				GitHub: GitHubConfig{
					Token: "test-token",
					Owner: "test-org",
					Repos: []string{},
				},
				Secrets: map[string]string{"SECRET1": "value1"},
			},
			wantErr: true,
			errMsg:  "at least one repository must be specified",
		},
		{
			name: "missing token",
			config: &Config{
				GitHub: GitHubConfig{
					Token: "",
					Owner: "test-org",
					Repos: []string{"repo1"},
				},
				Secrets: map[string]string{"SECRET1": "value1"},
			},
			wantErr: true,
			errMsg:  "github.token is required",
		},
		{
			name: "missing secrets",
			config: &Config{
				GitHub: GitHubConfig{
					Token: "test-token",
					Owner: "test-org",
					Repos: []string{"repo1"},
				},
				Secrets: map[string]string{},
			},
			wantErr: true,
			errMsg:  "at least one secret must be specified",
		},
		{
			name: "empty repo name",
			config: &Config{
				GitHub: GitHubConfig{
					Token: "test-token",
					Owner: "test-org",
					Repos: []string{""},
				},
				Secrets: map[string]string{"SECRET1": "value1"},
			},
			wantErr: true,
			errMsg:  "repository name cannot be empty",
		},
		{
			name: "empty secret key",
			config: &Config{
				GitHub: GitHubConfig{
					Token: "test-token",
					Owner: "test-org",
					Repos: []string{"repo1"},
				},
				Secrets: map[string]string{"": "value1"},
			},
			wantErr: true,
			errMsg:  "secret key cannot be empty",
		},
		{
			name: "empty secret value",
			config: &Config{
				GitHub: GitHubConfig{
					Token: "test-token",
					Owner: "test-org",
					Repos: []string{"repo1"},
				},
				Secrets: map[string]string{"SECRET1": ""},
			},
			wantErr: true,
			errMsg:  "secret value for 'SECRET1' cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConfig_ApplyOverrides(t *testing.T) {
	cfg := &Config{
		GitHub: GitHubConfig{
			Token: "original-token",
			Owner: "original-owner",
			Repos: []string{"repo1", "repo2"},
		},
		Secrets: map[string]string{"SECRET1": "value1"},
	}

	cfg.ApplyOverrides("new-token", "new-owner", []string{"repo3"})

	assert.Equal(t, "new-token", cfg.GitHub.Token)
	assert.Equal(t, "new-owner", cfg.GitHub.Owner)
	assert.Equal(t, []string{"repo3"}, cfg.GitHub.Repos)
}

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	tmpFile, err := os.CreateTemp("", "test-config-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	configContent := `
github:
  token: test-token
  owner: test-org
  repos:
    - repo1
    - repo2

secrets:
  SECRET1: "value1"
  SECRET2: "value2"
`

	_, err = tmpFile.WriteString(configContent)
	require.NoError(t, err)
	tmpFile.Close()

	cfg, err := LoadConfig(tmpFile.Name())
	require.NoError(t, err)

	assert.Equal(t, "test-token", cfg.GitHub.Token)
	assert.Equal(t, "test-org", cfg.GitHub.Owner)
	assert.Equal(t, []string{"repo1", "repo2"}, cfg.GitHub.Repos)
	assert.Equal(t, "value1", cfg.Secrets["SECRET1"])
	assert.Equal(t, "value2", cfg.Secrets["SECRET2"])
}

func TestLoadConfig_WithEnvToken(t *testing.T) {
	// Set environment variable
	os.Setenv(EnvTokenKey, "env-token")
	defer os.Unsetenv(EnvTokenKey)

	tmpFile, err := os.CreateTemp("", "test-config-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	configContent := `
github:
  owner: test-org
  repos:
    - repo1

secrets:
  SECRET1: "value1"
`

	_, err = tmpFile.WriteString(configContent)
	require.NoError(t, err)
	tmpFile.Close()

	cfg, err := LoadConfig(tmpFile.Name())
	require.NoError(t, err)

	assert.Equal(t, "env-token", cfg.GitHub.Token)
}

