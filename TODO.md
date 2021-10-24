# TODO

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
