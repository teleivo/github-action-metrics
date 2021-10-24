# TODO

* how do I get job ids and run ids?

first I need all runs

https://docs.github.com/en/rest/reference/actions#list-workflow-runs-for-a-repository

then for each run the jobs I am interested in

fetch the "jobs_url" inside the run
(https://docs.github.com/en/rest/reference/actions#get-a-job-for-a-workflow-run)

* get all runs. how to get them in a paced way. lets say via a GitHub action
  cron that runs twice a day and commits to this repo. Are there query params
  that allow me to filter by time in such a specific way? So that I do not miss
  any runs.
* autoformat code
* authenticate
