// Package cli provides command-line interface handlers.
package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/teleivo/github-action-metrics/internal/github"
	"github.com/teleivo/github-action-metrics/internal/storage"
)

// errFlagParse is a sentinel error indicating flag parsing failed.
// The flag package already printed the error, so main should not print again.
var errFlagParse = errors.New("flag parse error")

// FetchRunsConfig holds configuration for the fetch runs command.
type FetchRunsConfig struct {
	Repo        string
	Owner       string
	WorkflowID  int64
	Destination string
	Created     string
	WithJobs    bool
}

// FetchJobsConfig holds configuration for the fetch jobs command.
type FetchJobsConfig struct {
	Repo        string
	Owner       string
	WorkflowID  int64
	Destination string
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

// getGitHubToken returns the GitHub token from the GITHUB_TOKEN environment variable.
func getGitHubToken() string {
	return os.Getenv("GITHUB_TOKEN")
}

// HandleFetch handles the fetch command and its subcommands.
func HandleFetch(ctx context.Context, args []string, wErr io.Writer) (int, error) {
	if len(args) < 1 {
		printFetchUsage(wErr)
		return 2, nil
	}

	switch args[0] {
	case "runs":
		return handleFetchRuns(ctx, args[1:], wErr)
	case "jobs":
		return handleFetchJobs(ctx, args[1:], wErr)
	default:
		printFetchUsage(wErr)
		return 2, nil
	}
}

func printFetchUsage(w io.Writer) {
	_, _ = fmt.Fprintln(w, `Usage: gham fetch <command> [options]

Commands:
  runs    Fetch workflow runs from GitHub
  jobs    Fetch jobs for stored workflow runs

Run 'gham fetch <command> -h' for more information on a command.`)
}

func handleFetchRuns(ctx context.Context, args []string, wErr io.Writer) (int, error) {
	fs := flag.NewFlagSet("fetch runs", flag.ContinueOnError)
	fs.SetOutput(wErr)
	fs.Usage = func() {
		_, _ = fmt.Fprintln(wErr, `Usage: gham fetch runs [options]

Fetch latest GitHub action runs for a given workflow.

Requires GITHUB_TOKEN environment variable for authentication.

Options:`)
		fs.PrintDefaults()
	}

	repo := fs.String("repo", "", "GitHub repository (required)")
	owner := fs.String("owner", "", "Owner of GitHub repository (required)")
	workflowID := fs.Int64("workflow-id", 0, "Workflow ID of GitHub action (required)")
	destination := fs.String("destination", "", "Directory where payloads will be stored (required)")
	created := fs.String("created", "", "Date filter in format '2021-10-12' or '2021-10-29T22:40:19Z'")
	withJobs := fs.Bool("with-jobs", false, "Fetch jobs for fetched runs")

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0, nil
		}
		return 2, errFlagParse
	}

	// Validate required flags
	if *repo == "" || *owner == "" || *workflowID == 0 || *destination == "" {
		_, _ = fmt.Fprintln(wErr, "Error: -repo, -owner, -workflow-id, and -destination are required")
		fs.Usage()
		return 2, nil
	}

	dir, err := resolveDirectory(*destination)
	if err != nil {
		return 1, err
	}

	config := &FetchRunsConfig{
		Repo:        *repo,
		Owner:       *owner,
		WorkflowID:  *workflowID,
		Destination: dir,
		Created:     *created,
		WithJobs:    *withJobs,
	}

	if err := executeFetchRuns(ctx, config); err != nil {
		return 1, err
	}
	return 0, nil
}

func executeFetchRuns(ctx context.Context, config *FetchRunsConfig) error {
	store, err := storage.NewStore(config.Destination)
	if err != nil {
		return err
	}

	client := github.NewClient(getGitHubToken())

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

func handleFetchJobs(ctx context.Context, args []string, wErr io.Writer) (int, error) {
	fs := flag.NewFlagSet("fetch jobs", flag.ContinueOnError)
	fs.SetOutput(wErr)
	fs.Usage = func() {
		_, _ = fmt.Fprintln(wErr, `Usage: gham fetch jobs [options]

Fetch jobs for stored workflow runs.

Requires GITHUB_TOKEN environment variable for authentication.

Options:`)
		fs.PrintDefaults()
	}

	repo := fs.String("repo", "", "GitHub repository (required)")
	owner := fs.String("owner", "", "Owner of GitHub repository (required)")
	workflowID := fs.Int64("workflow-id", 0, "Workflow ID of GitHub action (required)")
	destination := fs.String("destination", "", "Directory where payloads are stored (required)")

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0, nil
		}
		return 2, errFlagParse
	}

	// Validate required flags
	if *repo == "" || *owner == "" || *workflowID == 0 || *destination == "" {
		_, _ = fmt.Fprintln(wErr, "Error: -repo, -owner, -workflow-id, and -destination are required")
		fs.Usage()
		return 2, nil
	}

	dir, err := resolveDirectory(*destination)
	if err != nil {
		return 1, err
	}

	config := &FetchJobsConfig{
		Repo:        *repo,
		Owner:       *owner,
		WorkflowID:  *workflowID,
		Destination: dir,
	}

	if err := executeFetchJobs(ctx, config); err != nil {
		return 1, err
	}
	return 0, nil
}

func executeFetchJobs(ctx context.Context, config *FetchJobsConfig) error {
	store, err := storage.NewStore(config.Destination)
	if err != nil {
		return err
	}

	client := github.NewClient(getGitHubToken())

	return github.FetchStoredRunJobs(ctx, client, config.Owner, config.Repo, config.WorkflowID, store)
}
