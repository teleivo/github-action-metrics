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

// HandleFetch handles the fetch command and its subcommands.
func HandleFetch(args []string) {
	if len(args) < 1 {
		printFetchUsage()
		os.Exit(1)
	}

	switch args[0] {
	case "runs":
		handleFetchRuns(args[1:])
	case "jobs":
		handleFetchJobs(args[1:])
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

func handleFetchRuns(args []string) {
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

	// Resolve and validate destination
	dir, err := filepath.Abs(*destination)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving destination: %v\n", err)
		os.Exit(1)
	}
	info, err := os.Stat(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "Error: %s must be a directory\n", dir)
		os.Exit(1)
	}

	// Get token from flag or environment
	authToken := *token
	if authToken == "" {
		authToken = os.Getenv("GITHUB_TOKEN")
	}

	config := &FetchRunsConfig{
		Repo:        *repo,
		Owner:       *owner,
		WorkflowID:  *workflowID,
		Destination: dir,
		Created:     *created,
		Token:       authToken,
		WithJobs:    *withJobs,
	}

	if err := executeFetchRuns(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func executeFetchRuns(config *FetchRunsConfig) error {
	store, err := storage.NewStore(config.Destination)
	if err != nil {
		return err
	}

	client := github.NewClient(config.Token)

	opts := &github.RunOptions{}
	if config.Created != "" {
		opts.Created = config.Created
	}

	ctx := context.Background()
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

func handleFetchJobs(args []string) {
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

	// Resolve and validate destination
	dir, err := filepath.Abs(*destination)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving destination: %v\n", err)
		os.Exit(1)
	}
	info, err := os.Stat(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "Error: %s must be a directory\n", dir)
		os.Exit(1)
	}

	// Get token from flag or environment
	authToken := *token
	if authToken == "" {
		authToken = os.Getenv("GITHUB_TOKEN")
	}

	config := &FetchJobsConfig{
		Repo:        *repo,
		Owner:       *owner,
		WorkflowID:  *workflowID,
		Destination: dir,
		Token:       authToken,
	}

	if err := executeFetchJobs(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func executeFetchJobs(config *FetchJobsConfig) error {
	store, err := storage.NewStore(config.Destination)
	if err != nil {
		return err
	}

	client := github.NewClient(config.Token)
	ctx := context.Background()

	return github.FetchStoredRunJobs(ctx, client, config.Owner, config.Repo, config.WorkflowID, store)
}
