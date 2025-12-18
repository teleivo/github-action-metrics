package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/teleivo/github-action-metrics/internal/storage"
)

// IndexRuns indexes workflow runs into Elasticsearch.
func IndexRuns(ctx context.Context, client *Client, store *storage.Store, workflowID int64) (*BulkResult, error) {
	docs := make(chan Document)

	go func() {
		defer close(docs)
		for data, err := range store.IterRuns(workflowID) {
			if err != nil {
				log.Printf("error reading run: %v", err)
				continue
			}

			var run map[string]any
			if err := json.Unmarshal(data, &run); err != nil {
				log.Printf("error unmarshaling run: %v", err)
				continue
			}

			runID, ok := run["id"].(float64)
			if !ok {
				log.Printf("run missing id field")
				continue
			}

			// Load jobs and compute duration
			jobsData, err := store.LoadJobs(workflowID, int64(runID))
			if err == nil {
				var jobs JobsResponse
				if err := json.Unmarshal(jobsData, &jobs); err == nil {
					if duration := ComputeRunDuration(&jobs); duration != nil {
						run["jobs_started_at"] = duration.JobsStartedAt
						run["jobs_started_at_id"] = duration.JobsStartedAtID
						run["jobs_started_at_name"] = duration.JobsStartedAtName
						run["jobs_started_at_url"] = duration.JobsStartedAtURL
						run["jobs_started_at_html_url"] = duration.JobsStartedAtHTMLURL
						run["jobs_completed_at"] = duration.JobsCompletedAt
						run["jobs_completed_at_id"] = duration.JobsCompletedAtID
						run["jobs_completed_at_name"] = duration.JobsCompletedAtName
						run["jobs_completed_at_url"] = duration.JobsCompletedAtURL
						run["jobs_completed_at_html_url"] = duration.JobsCompletedAtHTMLURL
					}
				}
			}

			docs <- Document{
				ID:   fmt.Sprintf("%.0f", runID),
				Body: run,
			}
		}
	}()

	result, err := client.BulkIndex(ctx, "runs", docs)
	if err != nil {
		return result, fmt.Errorf("indexing runs: %w", err)
	}

	log.Printf("indexed runs: %+v", result)
	return result, nil
}

// IndexJobs indexes workflow jobs into Elasticsearch.
func IndexJobs(ctx context.Context, client *Client, store *storage.Store, workflowID int64) (*BulkResult, error) {
	docs := make(chan Document)

	go func() {
		defer close(docs)
		for data, err := range store.IterJobs(workflowID) {
			if err != nil {
				log.Printf("error reading jobs: %v", err)
				continue
			}

			var jobsResp struct {
				Jobs []map[string]any `json:"jobs"`
			}
			if err := json.Unmarshal(data, &jobsResp); err != nil {
				log.Printf("error unmarshaling jobs: %v", err)
				continue
			}

			for _, job := range jobsResp.Jobs {
				jobID, ok := job["id"].(float64)
				if !ok {
					continue
				}
				docs <- Document{
					ID:   fmt.Sprintf("%.0f", jobID),
					Body: job,
				}
			}
		}
	}()

	result, err := client.BulkIndex(ctx, "jobs", docs)
	if err != nil {
		return result, fmt.Errorf("indexing jobs: %w", err)
	}

	log.Printf("indexed jobs: %+v", result)
	return result, nil
}

// IndexSteps indexes workflow steps into Elasticsearch.
func IndexSteps(ctx context.Context, client *Client, store *storage.Store, workflowID int64) (*BulkResult, error) {
	docs := make(chan Document)

	go func() {
		defer close(docs)
		for data, err := range store.IterJobs(workflowID) {
			if err != nil {
				log.Printf("error reading jobs: %v", err)
				continue
			}

			var jobsResp struct {
				Jobs []map[string]any `json:"jobs"`
			}
			if err := json.Unmarshal(data, &jobsResp); err != nil {
				log.Printf("error unmarshaling jobs: %v", err)
				continue
			}

			for _, job := range jobsResp.Jobs {
				jobID, _ := job["id"].(float64)
				jobName, _ := job["name"].(string)
				jobURL, _ := job["url"].(string)
				jobHTMLURL, _ := job["html_url"].(string)
				runID, _ := job["run_id"].(float64)
				runURL, _ := job["run_url"].(string)
				runAttempt, _ := job["run_attempt"].(float64)
				headSHA, _ := job["head_sha"].(string)

				stepsRaw, ok := job["steps"].([]any)
				if !ok {
					continue
				}

				for _, stepRaw := range stepsRaw {
					step, ok := stepRaw.(map[string]any)
					if !ok {
						continue
					}

					stepNumber, _ := step["number"].(float64)

					// Enrich step with job and run info
					step["job_id"] = jobID
					step["job_name"] = jobName
					step["job_url"] = jobURL
					step["job_html_url"] = jobHTMLURL
					step["run_id"] = runID
					step["run_url"] = runURL
					step["run_html_url"] = jobHTMLURL // Note: matches original TS behavior (bug)
					step["run_attempt"] = runAttempt
					step["head_sha"] = headSHA

					docs <- Document{
						ID:   fmt.Sprintf("%.0f", jobID) + "-" + strings.TrimSuffix(strconv.FormatFloat(stepNumber, 'f', -1, 64), ".0"),
						Body: step,
					}
				}
			}
		}
	}()

	result, err := client.BulkIndex(ctx, "steps", docs)
	if err != nil {
		return result, fmt.Errorf("indexing steps: %w", err)
	}

	log.Printf("indexed steps: %+v", result)
	return result, nil
}

// IndexAll indexes runs, jobs, and steps into Elasticsearch.
func IndexAll(ctx context.Context, client *Client, store *storage.Store, workflowID int64) error {
	if _, err := IndexRuns(ctx, client, store, workflowID); err != nil {
		return err
	}
	if _, err := IndexJobs(ctx, client, store, workflowID); err != nil {
		return err
	}
	if _, err := IndexSteps(ctx, client, store, workflowID); err != nil {
		return err
	}
	return nil
}
