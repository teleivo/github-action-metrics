// Package cli provides command-line interface handlers.
package cli

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/teleivo/github-action-metrics/internal/github"
	"github.com/teleivo/github-action-metrics/internal/storage"
)

// FetchRunsConfig holds configuration for the fetch runs command.
type FetchRunsConfig struct {
	Repo        string
	Owner       string
	WorkflowID  int64
	Destination string
	Created     string
	Token       string
	WithJobs    bool
}

// FetchJobsConfig holds configuration for the fetch jobs command.
type FetchJobsConfig struct {
	Repo        string
	Owner       string
	WorkflowID  int64
	Destination string
	Token       string
}

// resolveDirectory resolves a path to an absolute directory path.
// Returns an error if the path doesn't exist or isn't a directory.
func resolveDirectory(path string) (string, error) {
	dir, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("resolving path: %w", err)
	}
	info, err := os.Stat(dir)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("%s must be a directory", dir)
	}
	return dir, nil
}

// resolveToken returns the token from the flag or falls back to GITHUB_TOKEN env var.
func resolveToken(token string) string {
	if token != "" {
		return token
	}
	return os.Getenv("GITHUB_TOKEN")
}

// HandleFetch handles the fetch command and its subcommands.
func HandleFetch(ctx context.Context, args []string) {
	if len(args) < 1 {
		printFetchUsage()
		os.Exit(1)
	}

	switch args[0] {
	case "runs":
		handleFetchRuns(ctx, args[1:])
	case "jobs":
		handleFetchJobs(ctx, args[1:])
	default:
		printFetchUsage()
		os.Exit(1)
	}
}

func printFetchUsage() {
	fmt.Fprintln(os.Stderr, `Usage: gham fetch <command> [options]

Commands:
  runs    Fetch workflow runs from GitHub
  jobs    Fetch jobs for stored workflow runs

Run 'gham fetch <command> -h' for more information on a command.`)
}

func handleFetchRuns(ctx context.Context, args []string) {
	fs := flag.NewFlagSet("fetch runs", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, `Usage: gham fetch runs [options]

Fetch latest GitHub action runs for a given workflow.

Options:`)
		fs.PrintDefaults()
	}

	repo := fs.String("repo", "", "GitHub repository (required)")
	owner := fs.String("owner", "", "Owner of GitHub repository (required)")
	workflowID := fs.Int64("workflow-id", 0, "Workflow ID of GitHub action (required)")
	destination := fs.String("destination", "", "Directory where payloads will be stored (required)")
	created := fs.String("created", "", "Date filter in format '2021-10-12' or '2021-10-29T22:40:19Z'")
	token := fs.String("token", "", "GitHub access token (falls back to GITHUB_TOKEN env var)")
	withJobs := fs.Bool("with-jobs", false, "Fetch jobs for fetched runs")

	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	// Validate required flags
	if *repo == "" || *owner == "" || *workflowID == 0 || *destination == "" {
		fmt.Fprintln(os.Stderr, "Error: -repo, -owner, -workflow-id, and -destination are required")
		fs.Usage()
		os.Exit(1)
	}

	dir, err := resolveDirectory(*destination)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	config := &FetchRunsConfig{
		Repo:        *repo,
		Owner:       *owner,
		WorkflowID:  *workflowID,
		Destination: dir,
		Created:     *created,
		Token:       resolveToken(*token),
		WithJobs:    *withJobs,
	}

	if err := executeFetchRuns(ctx, config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func executeFetchRuns(ctx context.Context, config *FetchRunsConfig) error {
	store, err := storage.NewStore(config.Destination)
	if err != nil {
		return err
	}

	client := github.NewClient(config.Token)

	opts := &github.RunOptions{}
	if config.Created != "" {
		opts.Created = config.Created
	}

	runIDs, err := github.FetchRuns(ctx, client, config.Owner, config.Repo, config.WorkflowID, store, opts)
	if err != nil {
		return err
	}

	if config.WithJobs && len(runIDs) > 0 {
		if err := github.FetchJobs(ctx, client, config.Owner, config.Repo, config.WorkflowID, store, runIDs); err != nil {
			return err
		}
	}

	return nil
}

func handleFetchJobs(ctx context.Context, args []string) {
	fs := flag.NewFlagSet("fetch jobs", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, `Usage: gham fetch jobs [options]

Fetch jobs for stored workflow runs.

Options:`)
		fs.PrintDefaults()
	}

	repo := fs.String("repo", "", "GitHub repository (required)")
	owner := fs.String("owner", "", "Owner of GitHub repository (required)")
	workflowID := fs.Int64("workflow-id", 0, "Workflow ID of GitHub action (required)")
	destination := fs.String("destination", "", "Directory where payloads are stored (required)")
	token := fs.String("token", "", "GitHub access token (falls back to GITHUB_TOKEN env var)")

	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	// Validate required flags
	if *repo == "" || *owner == "" || *workflowID == 0 || *destination == "" {
		fmt.Fprintln(os.Stderr, "Error: -repo, -owner, -workflow-id, and -destination are required")
		fs.Usage()
		os.Exit(1)
	}

	dir, err := resolveDirectory(*destination)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	config := &FetchJobsConfig{
		Repo:        *repo,
		Owner:       *owner,
		WorkflowID:  *workflowID,
		Destination: dir,
		Token:       resolveToken(*token),
	}

	if err := executeFetchJobs(ctx, config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func executeFetchJobs(ctx context.Context, config *FetchJobsConfig) error {
	store, err := storage.NewStore(config.Destination)
	if err != nil {
		return err
	}

	client := github.NewClient(config.Token)

	return github.FetchStoredRunJobs(ctx, client, config.Owner, config.Repo, config.WorkflowID, store)
}
