package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
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
				slog.Warn("error reading run", "error", err)
				continue
			}

			var run map[string]any
			if err := json.Unmarshal(data, &run); err != nil {
				slog.Warn("error unmarshaling run", "error", err)
				continue
			}

			runID, ok := run["id"].(float64)
			if !ok {
				slog.Warn("run missing id field")
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
				ID:   strconv.FormatInt(int64(runID), 10),
				Body: run,
			}
		}
	}()

	result, err := client.BulkIndex(ctx, "runs", docs)
	if err != nil {
		return result, fmt.Errorf("indexing runs: %w", err)
	}

	slog.Info("indexed runs", "total", result.Total, "successful", result.Successful, "failed", result.Failed)
	return result, nil
}

// IndexJobs indexes workflow jobs into Elasticsearch.
func IndexJobs(ctx context.Context, client *Client, store *storage.Store, workflowID int64) (*BulkResult, error) {
	docs := make(chan Document)

	go func() {
		defer close(docs)
		for data, err := range store.IterJobs(workflowID) {
			if err != nil {
				slog.Warn("error reading jobs", "error", err)
				continue
			}

			var jobsResp struct {
				Jobs []map[string]any `json:"jobs"`
			}
			if err := json.Unmarshal(data, &jobsResp); err != nil {
				slog.Warn("error unmarshaling jobs", "error", err)
				continue
			}

			for _, job := range jobsResp.Jobs {
				jobID, ok := job["id"].(float64)
				if !ok {
					continue
				}
				docs <- Document{
					ID:   strconv.FormatInt(int64(jobID), 10),
					Body: job,
				}
			}
		}
	}()

	result, err := client.BulkIndex(ctx, "jobs", docs)
	if err != nil {
		return result, fmt.Errorf("indexing jobs: %w", err)
	}

	slog.Info("indexed jobs", "total", result.Total, "successful", result.Successful, "failed", result.Failed)
	return result, nil
}

// IndexSteps indexes workflow steps into Elasticsearch.
func IndexSteps(ctx context.Context, client *Client, store *storage.Store, workflowID int64) (*BulkResult, error) {
	docs := make(chan Document)

	go func() {
		defer close(docs)
		for data, err := range store.IterJobs(workflowID) {
			if err != nil {
				slog.Warn("error reading jobs", "error", err)
				continue
			}

			var jobsResp struct {
				Jobs []map[string]any `json:"jobs"`
			}
			if err := json.Unmarshal(data, &jobsResp); err != nil {
				slog.Warn("error unmarshaling jobs", "error", err)
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

				// Convert API URL to HTML URL: api.github.com/repos/... -> github.com/...
				runHTMLURL := strings.Replace(runURL, "api.github.com/repos", "github.com", 1)

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
					step["run_html_url"] = runHTMLURL
					step["run_attempt"] = runAttempt
					step["head_sha"] = headSHA

					docs <- Document{
						ID:   strconv.FormatInt(int64(jobID), 10) + "-" + strconv.FormatInt(int64(stepNumber), 10),
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

	slog.Info("indexed steps", "total", result.Total, "successful", result.Successful, "failed", result.Failed)
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
