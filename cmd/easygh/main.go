package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/spf13/cobra"

	"github.com/yourusername/easy_gh_secret/internal/cli"
	"github.com/yourusername/easy_gh_secret/internal/config"
	"github.com/yourusername/easy_gh_secret/internal/github"
	"github.com/yourusername/easy_gh_secret/internal/logger"
)

var (
	AZ_VERSION string = "1.0.0"
	AZ_UPDATE  string = "2025-12-10"
)

func main() {
	rootCmd := &cobra.Command{
		Use:          "easygh",
		Short:        "Manage GitHub Actions secrets across multiple repositories",
		Long:         `easygh is a CLI tool to manage GitHub Actions secrets across multiple repositories using a YAML configuration file.`,
		RunE:         run,
		SilenceUsage: true,
	}

	flags := &cli.Flags{}
	rootCmd.Flags().StringVarP(&flags.ConfigPath, "config", "c", "config.yaml", "Path to configuration file")
	rootCmd.Flags().StringVar(&flags.Token, "token", "", "GitHub token (overrides config file)")
	rootCmd.Flags().StringVar(&flags.Owner, "owner", "", "GitHub owner/organization (overrides config file)")
	rootCmd.Flags().StringVar(&flags.Repos, "repo", "", "Comma-separated list of repositories (overrides config file)")
	rootCmd.Flags().BoolVar(&flags.DryRun, "dry-run", false, "Show what would be done without making changes")
	rootCmd.Flags().BoolVar(&flags.ContinueOnError, "continue-on-error", false, "Continue processing other repositories on error")
	rootCmd.Flags().BoolVarP(&flags.Verbose, "verbose", "v", false, "Enable verbose logging")
	rootCmd.Flags().BoolVar(&flags.ShowVersion, "version", false, "Show version information")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	flags := &cli.Flags{}
	flags.ConfigPath, _ = cmd.Flags().GetString("config")
	flags.Token, _ = cmd.Flags().GetString("token")
	flags.Owner, _ = cmd.Flags().GetString("owner")
	flags.Repos, _ = cmd.Flags().GetString("repo")
	flags.DryRun, _ = cmd.Flags().GetBool("dry-run")
	flags.ContinueOnError, _ = cmd.Flags().GetBool("continue-on-error")
	flags.Verbose, _ = cmd.Flags().GetBool("verbose")
	flags.ShowVersion, _ = cmd.Flags().GetBool("version")

	// Initialize logger
	log := logger.New(flags.Verbose)

	// Show version if requested
	if flags.ShowVersion {
		fmt.Printf("Version: %s\n", AZ_VERSION)
		fmt.Printf("Build Time: %s\n", AZ_UPDATE)
		return nil
	}

	// Load configuration
	cfg, err := config.LoadConfigFromPath(flags.ConfigPath)
	if err != nil {
		log.Error("Failed to load configuration", "error", err)
		return err
	}

	// Apply CLI flag overrides
	repos := cli.ParseRepos(flags.Repos)
	cfg.ApplyOverrides(flags.Token, flags.Owner, repos)

	// Validate configuration again after overrides
	if err := cfg.Validate(); err != nil {
		log.Error("Configuration validation failed", "error", err)
		return err
	}

	// Create GitHub client
	ghClient := github.NewClient(cfg.GitHub.Token)

	// Execute the main logic
	ctx := context.Background()
	return execute(ctx, log, ghClient, cfg, flags)
}

func execute(ctx context.Context, log *logger.Logger, ghClient github.Client, cfg *config.Config, flags *cli.Flags) error {
	log.Info("Starting secrets management", "owner", cfg.GitHub.Owner, "repos", len(cfg.GitHub.Repos), "secrets", len(cfg.Secrets))

	if flags.DryRun {
		log.Info("DRY RUN MODE - No changes will be made")
	}

	// Create a cancellable context
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	var errorMutex sync.Mutex
	var errors []error

	// Process repositories concurrently
	for _, repo := range cfg.GitHub.Repos {
		wg.Add(1)
		go func(repoName string) {
			defer wg.Done()

			// Check if context is cancelled
			if ctx.Err() != nil {
				return
			}

			log.Info("Processing repository", "repo", repoName)

			repoErrors := processRepository(ctx, log, ghClient, cfg.GitHub.Owner, repoName, cfg.Secrets, flags.DryRun)

			if len(repoErrors) > 0 {
				errorMutex.Lock()
				errors = append(errors, repoErrors...)
				errorMutex.Unlock()

				if !flags.ContinueOnError {
					// Cancel context to stop other goroutines
					cancel()
				}
			}
		}(repo)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Report results
	if len(errors) > 0 {
		log.Error("Completed with errors", "error_count", len(errors))
		for _, err := range errors {
			log.Error("Error", "error", err)
		}
		return fmt.Errorf("failed with %d error(s)", len(errors))
	}

	log.Info("Successfully completed")
	return nil
}

func processRepository(ctx context.Context, log *logger.Logger, ghClient github.Client, owner, repo string, secrets map[string]string, dryRun bool) []error {
	var errors []error

	for secretName, secretValue := range secrets {
		// Check if context is cancelled
		if ctx.Err() != nil {
			return errors
		}

		if dryRun {
			// Try to get existing secret to show diff
			existingSecret, err := ghClient.GetSecret(ctx, owner, repo, secretName)
			if err != nil {
				log.Info("Would create secret", "repo", repo, "secret", secretName, "value", maskSecret(secretValue))
			} else {
				log.Info("Would update secret", "repo", repo, "secret", secretName, "existing", existingSecret.Name, "new_value", maskSecret(secretValue))
			}
		} else {
			if err := ghClient.SetSecret(ctx, owner, repo, secretName, secretValue); err != nil {
				log.Error("Failed to set secret", "repo", repo, "secret", secretName, "error", err)
				errors = append(errors, fmt.Errorf("repo %s/%s secret %s: %w", owner, repo, secretName, err))
				continue
			}
			log.Info("Successfully set secret", "repo", repo, "secret", secretName)
		}
	}

	return errors
}

func maskSecret(secret string) string {
	if len(secret) <= 4 {
		return "****"
	}
	return secret[:2] + "****" + secret[len(secret)-2:]
}

