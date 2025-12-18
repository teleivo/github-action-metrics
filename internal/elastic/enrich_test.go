package elastic

import (
	"testing"
)

func TestComputeRunDuration(t *testing.T) {
	tests := []struct {
		name string
		jobs *JobsResponse
		want *RunDuration
	}{
		{
			name: "two jobs returns first and last",
			jobs: &JobsResponse{
				TotalCount: 2,
				Jobs: []Job{
					{
						ID:          3865561946,
						Name:        "unit-test",
						URL:         "https://api.github.com/repos/dhis2/dhis2-core/actions/jobs/3865561946",
						HTMLURL:     "https://github.com/dhis2/dhis2-core/runs/3865561946?check_suite_focus=true",
						StartedAt:   "2021-10-12T01:56:28Z",
						CompletedAt: "2021-10-12T02:12:39Z",
					},
					{
						ID:          3865561982,
						Name:        "integration-test",
						URL:         "https://api.github.com/repos/dhis2/dhis2-core/actions/jobs/3865561982",
						HTMLURL:     "https://github.com/dhis2/dhis2-core/runs/3865561982?check_suite_focus=true",
						StartedAt:   "2021-10-12T01:57:25Z",
						CompletedAt: "2021-10-12T02:21:41Z",
					},
				},
			},
			want: &RunDuration{
				JobsStartedAt:          "2021-10-12T01:56:28Z",
				JobsStartedAtID:        3865561946,
				JobsStartedAtName:      "unit-test",
				JobsStartedAtURL:       "https://api.github.com/repos/dhis2/dhis2-core/actions/jobs/3865561946",
				JobsStartedAtHTMLURL:   "https://github.com/dhis2/dhis2-core/runs/3865561946?check_suite_focus=true",
				JobsCompletedAt:        "2021-10-12T02:21:41Z",
				JobsCompletedAtID:      3865561982,
				JobsCompletedAtName:    "integration-test",
				JobsCompletedAtURL:     "https://api.github.com/repos/dhis2/dhis2-core/actions/jobs/3865561982",
				JobsCompletedAtHTMLURL: "https://github.com/dhis2/dhis2-core/runs/3865561982?check_suite_focus=true",
			},
		},
		{
			name: "single job returns same job for start and end",
			jobs: &JobsResponse{
				TotalCount: 1,
				Jobs: []Job{
					{
						ID:          3865561946,
						Name:        "unit-test",
						URL:         "https://api.github.com/repos/dhis2/dhis2-core/actions/jobs/3865561946",
						HTMLURL:     "https://github.com/dhis2/dhis2-core/runs/3865561946?check_suite_focus=true",
						StartedAt:   "2021-10-12T01:57:28Z",
						CompletedAt: "2021-10-12T02:12:39Z",
					},
				},
			},
			want: &RunDuration{
				JobsStartedAt:          "2021-10-12T01:57:28Z",
				JobsStartedAtID:        3865561946,
				JobsStartedAtName:      "unit-test",
				JobsStartedAtURL:       "https://api.github.com/repos/dhis2/dhis2-core/actions/jobs/3865561946",
				JobsStartedAtHTMLURL:   "https://github.com/dhis2/dhis2-core/runs/3865561946?check_suite_focus=true",
				JobsCompletedAt:        "2021-10-12T02:12:39Z",
				JobsCompletedAtID:      3865561946,
				JobsCompletedAtName:    "unit-test",
				JobsCompletedAtURL:     "https://api.github.com/repos/dhis2/dhis2-core/actions/jobs/3865561946",
				JobsCompletedAtHTMLURL: "https://github.com/dhis2/dhis2-core/runs/3865561946?check_suite_focus=true",
			},
		},
		{
			name: "empty jobs returns nil",
			jobs: &JobsResponse{
				TotalCount: 0,
				Jobs:       []Job{},
			},
			want: nil,
		},
		{
			name: "nil jobs returns nil",
			jobs: nil,
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ComputeRunDuration(tt.jobs)
			if tt.want == nil {
				if got != nil {
					t.Errorf("ComputeRunDuration() = %+v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Errorf("ComputeRunDuration() = nil, want %+v", tt.want)
				return
			}
			if *got != *tt.want {
				t.Errorf("ComputeRunDuration() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
