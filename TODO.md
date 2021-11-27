# TODO

## CLI

- make fetching jobs after runs default? or configurable?
- cleanup CLI code (reuse options, extract action handlers?) reuse OctoKit and
  Elastic clients

## General

- reuse the OctoKit instance as I assume that has some state to make sure I
  keep below the ratelimit

- can I fetch multiple jobs at once? Currently every job leads to one request
  when I initially fetch all jobs I get (assume I hit the rate limit):

RequestError [HttpError]: request to https://api.github.com/repos/dhis2/dhis2-core/actions/runs/1282161011/jobs failed, reason: connect EHOSTUNREACH 140.82.114.6:443
at /home/ivo/code/dhis2/github-action-metrics/node_modules/@octokit/request/dist-node/index.js:108:11
at runMicrotasks (<anonymous>)
at processTicksAndRejections (internal/process/task_queues.js:97:5)
at async Job.doExecute (/home/ivo/code/dhis2/github-action-metrics/node_modules/bottleneck/light.js:405:18) {
status: 500,
request: {
method: 'GET',
url: 'https://api.github.com/repos/dhis2/dhis2-core/actions/runs/1282161011/jobs',
headers: {
accept: 'application/vnd.github.v3+json',
'user-agent': 'octokit-rest.js/1.7.0 octokit-core.js/3.5.1 Node.js/12.22.5 (linux; x64)',
authorization: 'token [REDACTED]'
},
request: {
hook: [Function: bound bound register],
retryCount: 3,
retries: 3,
retryAfter: 16
}
}
}

- rethink log messages
- add debug flag to reduce logging to a minimum

## Elastic

- create a backup/snapshot, also from kibana and try to restore?
- runs do not have the completed_at property. I think I would need to embedd
  jobs into runs so I could answer questions like `How long took successful runs of PRs to master in the last four weeks?` Or I could pre-process the
  jobs for a run and just add fields like created_at, completed_at for the
  first and last job to start/finish

## Whishlist

- can I fetch the workflowId so it can be passed by name in the CLI?
