// Package storage provides file-based JSON storage for GitHub Actions data.
package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"iter"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Store manages file-based storage for workflow runs and jobs.
type Store struct {
	baseDir string
}

// NewStore creates a new Store with the given base directory.
// Returns an error if the directory doesn't exist or isn't a directory.
func NewStore(baseDir string) (*Store, error) {
	info, err := os.Stat(baseDir)
	if err != nil {
		return nil, fmt.Errorf("invalid directory %q: %w", baseDir, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%q must be a directory", baseDir)
	}
	return &Store{baseDir: baseDir}, nil
}

// RunPath returns the file path for a workflow run.
func (s *Store) RunPath(workflowID, runID int64) string {
	return filepath.Join(s.baseDir, "workflows", fmt.Sprintf("%d", workflowID), "runs", fmt.Sprintf("%d.json", runID))
}

// JobPath returns the file path for jobs of a workflow run.
func (s *Store) JobPath(workflowID, runID int64) string {
	return filepath.Join(s.baseDir, "workflows", fmt.Sprintf("%d", workflowID), "jobs", fmt.Sprintf("%d.json", runID))
}

// RunExists checks if a run file exists.
func (s *Store) RunExists(workflowID, runID int64) bool {
	_, err := os.Stat(s.RunPath(workflowID, runID))
	return err == nil
}

// JobExists checks if a job file exists for the given run.
func (s *Store) JobExists(workflowID, runID int64) bool {
	_, err := os.Stat(s.JobPath(workflowID, runID))
	return err == nil
}

// SaveRun saves a workflow run as JSON.
func (s *Store) SaveRun(workflowID, runID int64, data json.RawMessage) error {
	path := s.RunPath(workflowID, runID)
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating directory %q: %w", dir, err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("writing run file %q: %w", path, err)
	}
	return nil
}

// SaveJobs saves workflow jobs as JSON.
func (s *Store) SaveJobs(workflowID, runID int64, data json.RawMessage) error {
	path := s.JobPath(workflowID, runID)
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating directory %q: %w", dir, err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("writing jobs file %q: %w", path, err)
	}
	return nil
}

// LoadRun loads a single run from storage.
func (s *Store) LoadRun(workflowID, runID int64) (json.RawMessage, error) {
	data, err := os.ReadFile(s.RunPath(workflowID, runID))
	if err != nil {
		return nil, fmt.Errorf("reading run file: %w", err)
	}
	return data, nil
}

// LoadJobs loads jobs for a single run from storage.
func (s *Store) LoadJobs(workflowID, runID int64) (json.RawMessage, error) {
	data, err := os.ReadFile(s.JobPath(workflowID, runID))
	if err != nil {
		return nil, fmt.Errorf("reading jobs file: %w", err)
	}
	return data, nil
}

// RunsDir returns the directory path for runs of a workflow.
func (s *Store) RunsDir(workflowID int64) string {
	return filepath.Join(s.baseDir, "workflows", fmt.Sprintf("%d", workflowID), "runs")
}

// JobsDir returns the directory path for jobs of a workflow.
func (s *Store) JobsDir(workflowID int64) string {
	return filepath.Join(s.baseDir, "workflows", fmt.Sprintf("%d", workflowID), "jobs")
}

// IterRuns iterates over all stored runs for a workflow.
func (s *Store) IterRuns(workflowID int64) iter.Seq2[json.RawMessage, error] {
	return func(yield func(json.RawMessage, error) bool) {
		dir := s.RunsDir(workflowID)
		entries, err := os.ReadDir(dir)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return
			}
			yield(nil, fmt.Errorf("reading runs directory %q: %w", dir, err))
			return
		}
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
				continue
			}
			data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
			if err != nil {
				if !yield(nil, fmt.Errorf("reading run file %q: %w", entry.Name(), err)) {
					return
				}
				continue
			}
			if !yield(data, nil) {
				return
			}
		}
	}
}

// IterJobs iterates over all stored job files for a workflow.
func (s *Store) IterJobs(workflowID int64) iter.Seq2[json.RawMessage, error] {
	return func(yield func(json.RawMessage, error) bool) {
		dir := s.JobsDir(workflowID)
		entries, err := os.ReadDir(dir)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return
			}
			yield(nil, fmt.Errorf("reading jobs directory %q: %w", dir, err))
			return
		}
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
				continue
			}
			data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
			if err != nil {
				if !yield(nil, fmt.Errorf("reading jobs file %q: %w", entry.Name(), err)) {
					return
				}
				continue
			}
			if !yield(data, nil) {
				return
			}
		}
	}
}

// ListStoredRunIDs returns all run IDs stored for a workflow.
func (s *Store) ListStoredRunIDs(workflowID int64) ([]int64, error) {
	dir := s.RunsDir(workflowID)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading runs directory %q: %w", dir, err)
	}

	var runIDs []int64
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		name := strings.TrimSuffix(entry.Name(), ".json")
		id, err := strconv.ParseInt(name, 10, 64)
		if err != nil {
			continue
		}
		runIDs = append(runIDs, id)
	}
	return runIDs, nil
}

// ListStoredRunIDsWithoutJobs returns run IDs that don't have corresponding job files.
func (s *Store) ListStoredRunIDsWithoutJobs(workflowID int64) ([]int64, error) {
	runIDs, err := s.ListStoredRunIDs(workflowID)
	if err != nil {
		return nil, err
	}

	var missing []int64
	for _, id := range runIDs {
		if !s.JobExists(workflowID, id) {
			missing = append(missing, id)
		}
	}
	return missing, nil
}
