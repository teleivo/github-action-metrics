# TODO

* how to get all runs in a paced way. lets say via a GitHub action
  cron that runs twice a day and commits to this repo. Are there query params
  that allow me to filter by time in such a specific way? So that I do not miss
  any runs.

Seems like there is an API I either did not see in the docs or its not listed
:)

Get the runs for a specific workflow

curl \
  -H "Accept: application/vnd.github.v3+json" \
  https://api.github.com/repos/dhis2/dhis2-core/actions/workflows/10954/runs

This reduces the number of runs.

via

https://docs.github.com/en/rest/reference/actions#list-workflow-runs-for-a-repository

limit myself to a few workflow_id's 

10954 - run-tests.yml
5591308 - run-api-tests.yml

then for each run fetch the "jobs_url" inside the run
(https://docs.github.com/en/rest/reference/actions#get-a-job-for-a-workflow-run)

* autoformat code
* create a small CLI out of it so I can easily get all runs for other workflows
* then get jobs for runs?
