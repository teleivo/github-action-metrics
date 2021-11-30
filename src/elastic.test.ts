import { runDuration } from "./elastic";

describe("runDuration", () => {
  test("should return timings of first and last job", () => {
    const jobs = {
      total_count: 2,
      jobs: [
        {
          id: 3865561946,
          run_id: 1331274863,
          run_url:
            "https://api.github.com/repos/dhis2/dhis2-core/actions/runs/1331274863",
          run_attempt: 1,
          node_id: "CR_kwDOA_1uaM7mZ8ta",
          head_sha: "d67c0556c038aede681a85f15d3bb65f7a11dcdd",
          url: "https://api.github.com/repos/dhis2/dhis2-core/actions/jobs/3865561946",
          html_url:
            "https://github.com/dhis2/dhis2-core/runs/3865561946?check_suite_focus=true",
          status: "completed",
          conclusion: "success",
          started_at: "2021-10-12T01:56:28Z",
          completed_at: "2021-10-12T02:12:39Z",
          name: "unit-test",
          steps: [
            {
              name: "Set up job",
              status: "completed",
              conclusion: "success",
              number: 1,
              started_at: "2021-10-12T01:56:28.000Z",
              completed_at: "2021-10-12T01:57:30.000Z",
            },
            {
              name: "Run actions/checkout@v2",
              status: "completed",
              conclusion: "success",
              number: 2,
              started_at: "2021-10-12T01:57:30.000Z",
              completed_at: "2021-10-12T01:57:34.000Z",
            },
            {
              name: "Set up JDK 11",
              status: "completed",
              conclusion: "success",
              number: 3,
              started_at: "2021-10-12T01:57:34.000Z",
              completed_at: "2021-10-12T01:57:39.000Z",
            },
            {
              name: "Cache maven artifacts",
              status: "completed",
              conclusion: "success",
              number: 4,
              started_at: "2021-10-12T01:57:39.000Z",
              completed_at: "2021-10-12T01:57:59.000Z",
            },
            {
              name: "Test core",
              status: "completed",
              conclusion: "success",
              number: 5,
              started_at: "2021-10-12T01:57:59.000Z",
              completed_at: "2021-10-12T02:11:12.000Z",
            },
            {
              name: "Test dhis-web",
              status: "completed",
              conclusion: "success",
              number: 6,
              started_at: "2021-10-12T02:11:12.000Z",
              completed_at: "2021-10-12T02:12:38.000Z",
            },
            {
              name: "Post Cache maven artifacts",
              status: "completed",
              conclusion: "success",
              number: 10,
              started_at: "2021-10-12T02:12:38.000Z",
              completed_at: "2021-10-12T02:12:39.000Z",
            },
            {
              name: "Post Set up JDK 11",
              status: "completed",
              conclusion: "success",
              number: 11,
              started_at: "2021-10-12T02:12:39.000Z",
              completed_at: "2021-10-12T02:12:39.000Z",
            },
            {
              name: "Post Run actions/checkout@v2",
              status: "completed",
              conclusion: "success",
              number: 12,
              started_at: "2021-10-12T02:12:39.000Z",
              completed_at: "2021-10-12T02:12:39.000Z",
            },
            {
              name: "Complete job",
              status: "completed",
              conclusion: "success",
              number: 13,
              started_at: "2021-10-12T02:12:39.000Z",
              completed_at: "2021-10-12T02:12:39.000Z",
            },
          ],
          check_run_url:
            "https://api.github.com/repos/dhis2/dhis2-core/check-runs/3865561946",
          labels: ["ubuntu-latest"],
          runner_id: 5,
          runner_name: "GitHub Actions 5",
          runner_group_id: 2,
          runner_group_name: "GitHub Actions",
        },
        {
          id: 3865561982,
          run_id: 1331274863,
          run_url:
            "https://api.github.com/repos/dhis2/dhis2-core/actions/runs/1331274863",
          run_attempt: 1,
          node_id: "CR_kwDOA_1uaM7mZ8t-",
          head_sha: "d67c0556c038aede681a85f15d3bb65f7a11dcdd",
          url: "https://api.github.com/repos/dhis2/dhis2-core/actions/jobs/3865561982",
          html_url:
            "https://github.com/dhis2/dhis2-core/runs/3865561982?check_suite_focus=true",
          status: "completed",
          conclusion: "success",
          started_at: "2021-10-12T01:57:25Z",
          completed_at: "2021-10-12T02:21:41Z",
          name: "integration-test",
          steps: [
            {
              name: "Set up job",
              status: "completed",
              conclusion: "success",
              number: 1,
              started_at: "2021-10-12T01:57:25.000Z",
              completed_at: "2021-10-12T01:57:28.000Z",
            },
            {
              name: "Run actions/checkout@v2",
              status: "completed",
              conclusion: "success",
              number: 2,
              started_at: "2021-10-12T01:57:28.000Z",
              completed_at: "2021-10-12T01:57:33.000Z",
            },
            {
              name: "Set up JDK 11",
              status: "completed",
              conclusion: "success",
              number: 3,
              started_at: "2021-10-12T01:57:33.000Z",
              completed_at: "2021-10-12T01:57:39.000Z",
            },
            {
              name: "Cache maven artifacts",
              status: "completed",
              conclusion: "success",
              number: 4,
              started_at: "2021-10-12T01:57:39.000Z",
              completed_at: "2021-10-12T01:57:59.000Z",
            },
            {
              name: "Run integration tests",
              status: "completed",
              conclusion: "success",
              number: 5,
              started_at: "2021-10-12T01:57:59.000Z",
              completed_at: "2021-10-12T02:21:41.000Z",
            },
            {
              name: "Post Cache maven artifacts",
              status: "completed",
              conclusion: "success",
              number: 8,
              started_at: "2021-10-12T02:21:41.000Z",
              completed_at: "2021-10-12T02:21:41.000Z",
            },
            {
              name: "Post Set up JDK 11",
              status: "completed",
              conclusion: "success",
              number: 9,
              started_at: "2021-10-12T02:21:41.000Z",
              completed_at: "2021-10-12T02:21:41.000Z",
            },
            {
              name: "Post Run actions/checkout@v2",
              status: "completed",
              conclusion: "success",
              number: 10,
              started_at: "2021-10-12T02:21:41.000Z",
              completed_at: "2021-10-12T02:21:41.000Z",
            },
            {
              name: "Complete job",
              status: "completed",
              conclusion: "success",
              number: 11,
              started_at: "2021-10-12T02:21:41.000Z",
              completed_at: "2021-10-12T02:21:41.000Z",
            },
          ],
          check_run_url:
            "https://api.github.com/repos/dhis2/dhis2-core/check-runs/3865561982",
          labels: ["ubuntu-latest"],
          runner_id: 4,
          runner_name: "GitHub Actions 4",
          runner_group_id: 2,
          runner_group_name: "GitHub Actions",
        },
      ],
    };

    expect(runDuration(jobs)).toStrictEqual({
      jobs_started_at: "2021-10-12T01:56:28Z",
      jobs_started_at_id: 3865561946,
      jobs_started_at_name: "unit-test",
      jobs_started_at_url:
        "https://api.github.com/repos/dhis2/dhis2-core/actions/jobs/3865561946",
      jobs_started_at_html_url:
        "https://github.com/dhis2/dhis2-core/runs/3865561946?check_suite_focus=true",
      jobs_completed_at: "2021-10-12T02:21:41Z",
      jobs_completed_at_id: 3865561982,
      jobs_completed_at_name: "integration-test",
      jobs_completed_at_url:
        "https://api.github.com/repos/dhis2/dhis2-core/actions/jobs/3865561982",
      jobs_completed_at_html_url:
        "https://github.com/dhis2/dhis2-core/runs/3865561982?check_suite_focus=true",
    });
  });

  test("should return timings of only job", () => {
    const jobs = {
      total_count: 1,
      jobs: [
        {
          id: 3865561946,
          run_id: 1331274863,
          run_url:
            "https://api.github.com/repos/dhis2/dhis2-core/actions/runs/1331274863",
          run_attempt: 1,
          node_id: "CR_kwDOA_1uaM7mZ8ta",
          head_sha: "d67c0556c038aede681a85f15d3bb65f7a11dcdd",
          url: "https://api.github.com/repos/dhis2/dhis2-core/actions/jobs/3865561946",
          html_url:
            "https://github.com/dhis2/dhis2-core/runs/3865561946?check_suite_focus=true",
          status: "completed",
          conclusion: "success",
          started_at: "2021-10-12T01:57:28Z",
          completed_at: "2021-10-12T02:12:39Z",
          name: "unit-test",
          steps: [
            {
              name: "Set up job",
              status: "completed",
              conclusion: "success",
              number: 1,
              started_at: "2021-10-12T01:57:28.000Z",
              completed_at: "2021-10-12T01:57:30.000Z",
            },
            {
              name: "Run actions/checkout@v2",
              status: "completed",
              conclusion: "success",
              number: 2,
              started_at: "2021-10-12T01:57:30.000Z",
              completed_at: "2021-10-12T01:57:34.000Z",
            },
            {
              name: "Set up JDK 11",
              status: "completed",
              conclusion: "success",
              number: 3,
              started_at: "2021-10-12T01:57:34.000Z",
              completed_at: "2021-10-12T01:57:39.000Z",
            },
            {
              name: "Cache maven artifacts",
              status: "completed",
              conclusion: "success",
              number: 4,
              started_at: "2021-10-12T01:57:39.000Z",
              completed_at: "2021-10-12T01:57:59.000Z",
            },
            {
              name: "Test core",
              status: "completed",
              conclusion: "success",
              number: 5,
              started_at: "2021-10-12T01:57:59.000Z",
              completed_at: "2021-10-12T02:11:12.000Z",
            },
            {
              name: "Test dhis-web",
              status: "completed",
              conclusion: "success",
              number: 6,
              started_at: "2021-10-12T02:11:12.000Z",
              completed_at: "2021-10-12T02:12:38.000Z",
            },
            {
              name: "Post Cache maven artifacts",
              status: "completed",
              conclusion: "success",
              number: 10,
              started_at: "2021-10-12T02:12:38.000Z",
              completed_at: "2021-10-12T02:12:39.000Z",
            },
            {
              name: "Post Set up JDK 11",
              status: "completed",
              conclusion: "success",
              number: 11,
              started_at: "2021-10-12T02:12:39.000Z",
              completed_at: "2021-10-12T02:12:39.000Z",
            },
            {
              name: "Post Run actions/checkout@v2",
              status: "completed",
              conclusion: "success",
              number: 12,
              started_at: "2021-10-12T02:12:39.000Z",
              completed_at: "2021-10-12T02:12:39.000Z",
            },
            {
              name: "Complete job",
              status: "completed",
              conclusion: "success",
              number: 13,
              started_at: "2021-10-12T02:12:39.000Z",
              completed_at: "2021-10-12T02:12:39.000Z",
            },
          ],
          check_run_url:
            "https://api.github.com/repos/dhis2/dhis2-core/check-runs/3865561946",
          labels: ["ubuntu-latest"],
          runner_id: 5,
          runner_name: "GitHub Actions 5",
          runner_group_id: 2,
          runner_group_name: "GitHub Actions",
        },
      ],
    };

    expect(runDuration(jobs)).toStrictEqual({
      jobs_started_at: "2021-10-12T01:57:28Z",
      jobs_started_at_id: 3865561946,
      jobs_started_at_name: "unit-test",
      jobs_started_at_url:
        "https://api.github.com/repos/dhis2/dhis2-core/actions/jobs/3865561946",
      jobs_started_at_html_url:
        "https://github.com/dhis2/dhis2-core/runs/3865561946?check_suite_focus=true",
      jobs_completed_at: "2021-10-12T02:12:39Z",
      jobs_completed_at_id: 3865561946,
      jobs_completed_at_name: "unit-test",
      jobs_completed_at_url:
        "https://api.github.com/repos/dhis2/dhis2-core/actions/jobs/3865561946",
      jobs_completed_at_html_url:
        "https://github.com/dhis2/dhis2-core/runs/3865561946?check_suite_focus=true",
    });
  });

  test("should return empty if no jobs are given", () => {
    const jobs = {
      total_count: 0,
      jobs: [],
    };

    expect(runDuration(jobs)).toStrictEqual({});
  });
});
