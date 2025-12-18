package github

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/go-github/v67/github"
	"github.com/teleivo/github-action-metrics/internal/storage"
)

// FetchJobs fetches jobs for the given run IDs and stores them.
func FetchJobs(ctx context.Context, client *Client, owner, repo string, workflowID int64, store *storage.Store, runIDs []int64) error {
	for _, runID := range runIDs {
		if err := fetchJobsForRun(ctx, client, owner, repo, workflowID, store, runID); err != nil {
			log.Printf("failed to fetch jobs for run #%d: %v", runID, err)
		}
	}
	return nil
}

// FetchStoredRunJobs fetches jobs for all stored runs that don't have jobs yet.
func FetchStoredRunJobs(ctx context.Context, client *Client, owner, repo string, workflowID int64, store *storage.Store) error {
	runIDs, err := store.ListStoredRunIDsWithoutJobs(workflowID)
	if err != nil {
		return fmt.Errorf("listing runs without jobs: %w", err)
	}

	if len(runIDs) == 0 {
		log.Println("no runs without jobs found")
		return nil
	}

	log.Printf("fetching jobs for %d runs", len(runIDs))
	return FetchJobs(ctx, client, owner, repo, workflowID, store, runIDs)
}

func fetchJobsForRun(ctx context.Context, client *Client, owner, repo string, workflowID int64, store *storage.Store, runID int64) error {
	log.Printf("fetching jobs for run #%d", runID)

	opts := &github.ListWorkflowJobsOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	var allJobs []*github.WorkflowJob

	for {
		jobs, resp, err := client.Actions().ListWorkflowJobs(ctx, owner, repo, runID, opts)
		if err != nil {
			return fmt.Errorf("listing jobs for run #%d: %w", runID, err)
		}

		allJobs = append(allJobs, jobs.Jobs...)

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	log.Printf("fetched %d jobs for run #%d", len(allJobs), runID)

	// Store in same format as the TypeScript version: {"jobs": [...]}
	jobsResponse := struct {
		Jobs []*github.WorkflowJob `json:"jobs"`
	}{
		Jobs: allJobs,
	}

	data, err := json.Marshal(jobsResponse)
	if err != nil {
		return fmt.Errorf("marshaling jobs: %w", err)
	}

	if err := store.SaveJobs(workflowID, runID, data); err != nil {
		return fmt.Errorf("saving jobs: %w", err)
	}

	return nil
}
