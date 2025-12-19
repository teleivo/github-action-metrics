package github

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/go-github/v67/github"
	"github.com/teleivo/github-action-metrics/internal/storage"
)

// RunOptions configures the FetchRuns operation.
type RunOptions struct {
	Created string // Date filter in format '2021-10-12' or '2021-10-29T22:40:19Z'
}

// FetchRuns fetches workflow runs from GitHub and stores them locally.
// It skips runs that already exist in storage.
// Returns the IDs of newly fetched runs.
func FetchRuns(ctx context.Context, client *Client, owner, repo string, workflowID int64, store *storage.Store, opts *RunOptions) ([]int64, error) {
	listOpts := &github.ListWorkflowRunsOptions{
		Event:  "pull_request",
		Status: "completed",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	if opts != nil && opts.Created != "" {
		listOpts.Created = opts.Created
	}

	var fetchedRunIDs []int64

	for {
		runs, resp, err := client.Actions().ListWorkflowRunsByID(ctx, owner, repo, workflowID, listOpts)
		if err != nil {
			return fetchedRunIDs, fmt.Errorf("listing workflow runs: %w", err)
		}

		for _, run := range runs.WorkflowRuns {
			runID := run.GetID()
			slog.Debug("processing run", "run_id", runID)

			if store.RunExists(workflowID, runID) {
				slog.Debug("run already exists", "run_id", runID, "path", store.RunPath(workflowID, runID))
				continue
			}

			data, err := json.Marshal(run)
			if err != nil {
				slog.Warn("failed to marshal run", "run_id", runID, "error", err)
				continue
			}

			if err := store.SaveRun(workflowID, runID, data); err != nil {
				slog.Warn("failed to save run", "run_id", runID, "error", err)
				continue
			}

			fetchedRunIDs = append(fetchedRunIDs, runID)
		}

		if resp.NextPage == 0 {
			break
		}
		listOpts.Page = resp.NextPage
	}

	return fetchedRunIDs, nil
}
