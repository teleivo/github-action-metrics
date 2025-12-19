package cli

import (
	"context"
	"flag"
	"fmt"
	"io"

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
func HandleIndex(ctx context.Context, args []string, wErr io.Writer) (int, error) {
	if len(args) < 1 {
		printIndexUsage(wErr)
		return 2, nil
	}

	switch args[0] {
	case "runs":
		return handleIndexRuns(ctx, args[1:], wErr)
	case "jobs":
		return handleIndexJobs(ctx, args[1:], wErr)
	case "steps":
		return handleIndexSteps(ctx, args[1:], wErr)
	case "all":
		return handleIndexAll(ctx, args[1:], wErr)
	default:
		printIndexUsage(wErr)
		return 2, nil
	}
}

func printIndexUsage(w io.Writer) {
	fmt.Fprintln(w, `Usage: gham index <command> [options]

Commands:
  runs    Index workflow runs in Elasticsearch
  jobs    Index workflow jobs in Elasticsearch
  steps   Index workflow steps in Elasticsearch
  all     Index runs, jobs, and steps in Elasticsearch

Run 'gham index <command> -h' for more information on a command.`)
}

func parseIndexFlags(name string, args []string, wErr io.Writer) (*IndexConfig, int, error) {
	fs := flag.NewFlagSet("index "+name, flag.ContinueOnError)
	fs.SetOutput(wErr)
	fs.Usage = func() {
		fmt.Fprintf(wErr, "Usage: gham index %s [options]\n\nIndex workflow %s in Elasticsearch.\n\nOptions:\n", name, name)
		fs.PrintDefaults()
	}

	url := fs.String("url", "", "Elasticsearch URL (required)")
	workflowID := fs.Int64("workflow-id", 0, "Workflow ID of GitHub action (required)")
	source := fs.String("source", "", "Directory where GitHub action payloads are stored (required)")
	user := fs.String("user", "", "Elasticsearch basic authentication user (required)")
	password := fs.String("password", "", "Elasticsearch basic authentication password (required)")

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return nil, 0, nil
		}
		return nil, 2, errFlagParse
	}

	// Validate required flags
	if *url == "" || *workflowID == 0 || *source == "" || *user == "" || *password == "" {
		fmt.Fprintln(wErr, "Error: -url, -workflow-id, -source, -user, and -password are required")
		fs.Usage()
		return nil, 2, nil
	}

	dir, err := resolveDirectory(*source)
	if err != nil {
		return nil, 1, err
	}

	return &IndexConfig{
		URL:        *url,
		WorkflowID: *workflowID,
		Source:     dir,
		User:       *user,
		Password:   *password,
	}, 0, nil
}

func handleIndexRuns(ctx context.Context, args []string, wErr io.Writer) (int, error) {
	config, code, err := parseIndexFlags("runs", args, wErr)
	if config == nil {
		return code, err
	}

	store, err := storage.NewStore(config.Source)
	if err != nil {
		return 1, err
	}

	client := elastic.NewClient(config.URL, config.User, config.Password)

	if _, err := elastic.IndexRuns(ctx, client, store, config.WorkflowID); err != nil {
		return 1, err
	}
	return 0, nil
}

func handleIndexJobs(ctx context.Context, args []string, wErr io.Writer) (int, error) {
	config, code, err := parseIndexFlags("jobs", args, wErr)
	if config == nil {
		return code, err
	}

	store, err := storage.NewStore(config.Source)
	if err != nil {
		return 1, err
	}

	client := elastic.NewClient(config.URL, config.User, config.Password)

	if _, err := elastic.IndexJobs(ctx, client, store, config.WorkflowID); err != nil {
		return 1, err
	}
	return 0, nil
}

func handleIndexSteps(ctx context.Context, args []string, wErr io.Writer) (int, error) {
	config, code, err := parseIndexFlags("steps", args, wErr)
	if config == nil {
		return code, err
	}

	store, err := storage.NewStore(config.Source)
	if err != nil {
		return 1, err
	}

	client := elastic.NewClient(config.URL, config.User, config.Password)

	if _, err := elastic.IndexSteps(ctx, client, store, config.WorkflowID); err != nil {
		return 1, err
	}
	return 0, nil
}

func handleIndexAll(ctx context.Context, args []string, wErr io.Writer) (int, error) {
	config, code, err := parseIndexFlags("all", args, wErr)
	if config == nil {
		return code, err
	}

	store, err := storage.NewStore(config.Source)
	if err != nil {
		return 1, err
	}

	client := elastic.NewClient(config.URL, config.User, config.Password)

	if err := elastic.IndexAll(ctx, client, store, config.WorkflowID); err != nil {
		return 1, err
	}
	return 0, nil
}
