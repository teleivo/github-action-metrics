# TODO

- change released version to be 0.0.1
- add a changelog to the release
- add a pre-commit hook for checking semantic commits so release works as
  expected
- create README. mention without token you risk hitting rate limit

## CLI

- make fetching jobs after runs default? or configurable?
- cleanup CLI code (reuse options, extract action handlers?)

## GitHub actions

- create GitHub workflow to create a release once a version tag is pushed
- publish CLI to GitHub packages
- use CLI package in other repo in a GitHub action using cron to fetch and
  store the action payloads
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

- setup ELK stack
- ingest data
- create first visualization for for example how did the "Run integration tests" step evolve over time?

## Elastic

- automate kibana index pattern creation?
- and runtime field `duration` format into human readable
- can I automatically index steps array separately? Or at least not let ES
  merge the array into one single field?
- get rid of security warning in kibana
  https://www.elastic.co/guide/en/elasticsearch/reference/7.15/security-minimal-setup.html

## Whishlist

- can I fetch the workflowId so it can be passed by name in the CLI?
