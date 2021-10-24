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

* what if the total number changes? how to first get the total number, check
  how many I have already stored and only fetch whats left?
* get job for each run

limit myself to a few workflow_id's 

10954 - run-tests.yml
5591308 - run-api-tests.yml

then for each run fetch the "jobs_url" inside the run
(https://docs.github.com/en/rest/reference/actions#get-a-job-for-a-workflow-run)

* if I start with fetchAllRuns() and store the last page in data/workflow/10954/runs/lastPage
I can resume using the lastPage just in case the lastPage was not full

* autoformat code
* create a small CLI out of it so I can easily get all runs for other workflows
* then get jobs for runs?
* setup ELK stack
* ingest data
* create first visualization for for example how did the "Run integration tests" step evolve over time?
