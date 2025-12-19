package elastic

import (
	"time"
)

// RunDuration contains timing information about the first and last jobs in a run.
type RunDuration struct {
	JobsStartedAt        string `json:"jobs_started_at"`
	JobsStartedAtID      int64  `json:"jobs_started_at_id"`
	JobsStartedAtName    string `json:"jobs_started_at_name"`
	JobsStartedAtURL     string `json:"jobs_started_at_url"`
	JobsStartedAtHTMLURL string `json:"jobs_started_at_html_url"`

	JobsCompletedAt        string `json:"jobs_completed_at"`
	JobsCompletedAtID      int64  `json:"jobs_completed_at_id"`
	JobsCompletedAtName    string `json:"jobs_completed_at_name"`
	JobsCompletedAtURL     string `json:"jobs_completed_at_url"`
	JobsCompletedAtHTMLURL string `json:"jobs_completed_at_html_url"`
}

// Job represents a workflow job with timing information.
type Job struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	HTMLURL     string `json:"html_url"`
	StartedAt   string `json:"started_at"`
	CompletedAt string `json:"completed_at"`
}

// JobsResponse represents the response from GitHub's jobs API.
type JobsResponse struct {
	TotalCount int   `json:"total_count"`
	Jobs       []Job `json:"jobs"`
}

// ComputeRunDuration calculates the earliest start time and latest completion time
// across all jobs in a run. Returns nil if there are no jobs.
func ComputeRunDuration(jobs *JobsResponse) *RunDuration {
	if jobs == nil || len(jobs.Jobs) == 0 {
		return nil
	}

	var result RunDuration
	var earliestStart, latestComplete time.Time

	for _, job := range jobs.Jobs {
		startedAt, err := time.Parse(time.RFC3339, job.StartedAt)
		if err != nil {
			continue
		}
		completedAt, err := time.Parse(time.RFC3339, job.CompletedAt)
		if err != nil {
			continue
		}

		if earliestStart.IsZero() || startedAt.Before(earliestStart) {
			earliestStart = startedAt
			result.JobsStartedAt = job.StartedAt
			result.JobsStartedAtID = job.ID
			result.JobsStartedAtName = job.Name
			result.JobsStartedAtURL = job.URL
			result.JobsStartedAtHTMLURL = job.HTMLURL
		}

		if latestComplete.IsZero() || completedAt.After(latestComplete) {
			latestComplete = completedAt
			result.JobsCompletedAt = job.CompletedAt
			result.JobsCompletedAtID = job.ID
			result.JobsCompletedAtName = job.Name
			result.JobsCompletedAtURL = job.URL
			result.JobsCompletedAtHTMLURL = job.HTMLURL
		}
	}

	return &result
}
