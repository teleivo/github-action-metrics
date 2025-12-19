package cli

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/teleivo/github-action-metrics/internal/elastic"
	"github.com/teleivo/github-action-metrics/internal/storage"
)

// IndexConfig holds configuration for index commands.
type IndexConfig struct {
	URL        string
	WorkflowID int64
	Source     string
	User       string
	Password   string
}

// HandleIndex handles the index command and its subcommands.
func HandleIndex(ctx context.Context, args []string) {
	if len(args) < 1 {
		printIndexUsage()
		os.Exit(1)
	}

	switch args[0] {
	case "runs":
		handleIndexRuns(ctx, args[1:])
	case "jobs":
		handleIndexJobs(ctx, args[1:])
	case "steps":
		handleIndexSteps(ctx, args[1:])
	case "all":
		handleIndexAll(ctx, args[1:])
	default:
		printIndexUsage()
		os.Exit(1)
	}
}

func printIndexUsage() {
	fmt.Fprintln(os.Stderr, `Usage: gham index <command> [options]

Commands:
  runs    Index workflow runs in Elasticsearch
  jobs    Index workflow jobs in Elasticsearch
  steps   Index workflow steps in Elasticsearch
  all     Index runs, jobs, and steps in Elasticsearch

Run 'gham index <command> -h' for more information on a command.`)
}

func parseIndexFlags(name string, args []string) *IndexConfig {
	fs := flag.NewFlagSet("index "+name, flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: gham index %s [options]\n\nIndex workflow %s in Elasticsearch.\n\nOptions:\n", name, name)
		fs.PrintDefaults()
	}

	url := fs.String("url", "", "Elasticsearch URL (required)")
	workflowID := fs.Int64("workflow-id", 0, "Workflow ID of GitHub action (required)")
	source := fs.String("source", "", "Directory where GitHub action payloads are stored (required)")
	user := fs.String("user", "", "Elasticsearch basic authentication user (required)")
	password := fs.String("password", "", "Elasticsearch basic authentication password (required)")

	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	// Validate required flags
	if *url == "" || *workflowID == 0 || *source == "" || *user == "" || *password == "" {
		fmt.Fprintln(os.Stderr, "Error: -url, -workflow-id, -source, -user, and -password are required")
		fs.Usage()
		os.Exit(1)
	}

	dir, err := resolveDirectory(*source)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	return &IndexConfig{
		URL:        *url,
		WorkflowID: *workflowID,
		Source:     dir,
		User:       *user,
		Password:   *password,
	}
}

func handleIndexRuns(ctx context.Context, args []string) {
	config := parseIndexFlags("runs", args)

	store, err := storage.NewStore(config.Source)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	client := elastic.NewClient(config.URL, config.User, config.Password)

	if _, err := elastic.IndexRuns(ctx, client, store, config.WorkflowID); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func handleIndexJobs(ctx context.Context, args []string) {
	config := parseIndexFlags("jobs", args)

	store, err := storage.NewStore(config.Source)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	client := elastic.NewClient(config.URL, config.User, config.Password)

	if _, err := elastic.IndexJobs(ctx, client, store, config.WorkflowID); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func handleIndexSteps(ctx context.Context, args []string) {
	config := parseIndexFlags("steps", args)

	store, err := storage.NewStore(config.Source)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	client := elastic.NewClient(config.URL, config.User, config.Password)

	if _, err := elastic.IndexSteps(ctx, client, store, config.WorkflowID); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func handleIndexAll(ctx context.Context, args []string) {
	config := parseIndexFlags("all", args)

	store, err := storage.NewStore(config.Source)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	client := elastic.NewClient(config.URL, config.User, config.Password)

	if err := elastic.IndexAll(ctx, client, store, config.WorkflowID); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
