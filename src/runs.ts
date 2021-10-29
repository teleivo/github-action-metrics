import * as fs from "fs";
import * as path from "path";

import { Octokit } from "octokit";
import { RestEndpointMethodTypes } from "@octokit/plugin-rest-endpoint-methods";

// Although not specified in the docs https://docs.github.com/en/rest/reference/actions#list-workflow-runs
// the API only returns the 10 last pages. So we can get the last ~1000 runs
// with per_page set to 100. Older ones can be fetched using the parameter
// created https://github.blog/changelog/2021-09-02-github-actions-filter-workflow-runs-by-created-date/
// You can thus run this every day and it will store the latest runs you have
// not yet stored locally.
async function fetchLatest(
  params: RestEndpointMethodTypes["actions"]["listWorkflowRuns"]["parameters"],
  directory: string,
  token: string
) {
  const octokit = new Octokit({
    auth: token,
  });

  octokit.hook.after("request", async (response, options) => {
    const ratelimit = `${response.headers["x-ratelimit-used"]}/${response.headers["x-ratelimit-limit"]}`;
    console.log(
      `requested ${options.method} ${options.url}: ${response.status} ratelimit ${ratelimit}`
    );
  });
  octokit.hook.error("request", async (error, _) => {
    console.error(error);
  });

  try {
    const iterator = octokit.paginate.iterator(
      octokit.rest.actions.listWorkflowRuns,
      {
        ...params,
        event: "pull_request",
        status: "completed",
        per_page: 100,
      }
    );

    for await (const { data: runs } of iterator) {
      for (const run of runs) {
        console.log("Run #%d", run.id);
        const file = path.join(
          directory,
          `/workflows/${params.workflow_id}/runs/${run.id}.json`
        );
        if (fs.existsSync(file)) {
          console.log("run #%d already exists in %s", run.id, file);
          continue;
        }
        const fd = fs.openSync(file, "w");
        fs.writeFileSync(fd, JSON.stringify(run));
        // TODO ensure I always close the file
        fs.closeSync(fd);
      }
    }
  } catch (error) {
    console.error(error);
    return;
  }
}

async function fetchCreatedOn(
  params: RestEndpointMethodTypes["actions"]["listWorkflowRuns"]["parameters"],
  directory: string,
  created: string,
  token: string
) {
  const octokit = new Octokit({
    auth: token,
  });

  octokit.hook.after("request", async (response, options) => {
    const ratelimit = `${response.headers["x-ratelimit-used"]}/${response.headers["x-ratelimit-limit"]}`;
    console.log(
      `requested ${options.method} ${options.url}: ${response.status} ratelimit ${ratelimit}`
    );
  });
  octokit.hook.error("request", async (error, _) => {
    console.error(error);
  });

  try {
    const response = await octokit.rest.actions.listWorkflowRuns({
      ...params,
      event: "pull_request",
      status: "completed",
      created,
    });
    if (response.status != 200) {
      console.log(
        `failed to fetch runs, API responded with status ${response.status}`
      );
      return;
    }

    for (const run of response.data.workflow_runs) {
      console.log("Run #%d", run.id);
      const file = path.join(
        directory,
        `/workflows/${params.workflow_id}/runs/${run.id}.json`
      );
      if (fs.existsSync(file)) {
        console.log("run #%d already exists in %s", run.id, file);
        continue;
      }
      const fd = fs.openSync(file, "w");
      fs.writeFileSync(fd, JSON.stringify(run));
      // TODO ensure I always close the file
      fs.closeSync(fd);
    }
  } catch (error) {
    console.error(error);
    return;
  }
}

export function fetchLatestRuns(
  repo: string,
  owner: string,
  workflowId: number,
  directory: string,
  token: string
) {
  fetchLatest(
    {
      repo,
      owner,
      workflow_id: workflowId,
    },
    directory,
    token
  );
}

export function fetchRunsCreatedOn(
  repo: string,
  owner: string,
  workflowId: number,
  directory: string,
  created: string,
  token: string
) {
  fetchCreatedOn(
    {
      repo,
      owner,
      workflow_id: workflowId,
    },
    directory,
    created,
    token
  );
}
