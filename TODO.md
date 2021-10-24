# TODO

* how to get all runs in a paced way. lets say via a GitHub action
  cron that runs twice a day and commits to this repo. Are there query params
  that allow me to filter by time in such a specific way? So that I do not miss
  any runs.

via

https://docs.github.com/en/rest/reference/actions#list-workflow-runs-for-a-repository

then for each run fetch the "jobs_url" inside the run
(https://docs.github.com/en/rest/reference/actions#get-a-job-for-a-workflow-run)

* autoformat code
