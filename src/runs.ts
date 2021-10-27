import { Octokit } from "octokit";
import { openSync, closeSync, existsSync, writeFileSync } from "fs";
import { RestEndpointMethodTypes } from "@octokit/plugin-rest-endpoint-methods";

const octokit = new Octokit({
  auth: process.env.GITHUB_TOKEN,
});

// Although not specified in the docs https://docs.github.com/en/rest/reference/actions#list-workflow-runs
// the API only returns the 10 last pages. So we can get the last ~1000 runs
// with per_page set to 100. Older ones can be fetched using the parameter
// created https://github.blog/changelog/2021-09-02-github-actions-filter-workflow-runs-by-created-date/
// You can thus run this every day and it will store the latest runs you have
// not yet stored locally.
async function fetchLatestRuns(
  params: RestEndpointMethodTypes["actions"]["listWorkflowRuns"]["parameters"]
) {
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
        const file = `data/workflows/${params.workflow_id}/runs/${run.id}.json`;
        if (existsSync(file)) {
          console.log("run #%d already exists in %s", run.id, file);
          continue;
        }
        const fd = openSync(file, "w");
        writeFileSync(fd, JSON.stringify(run));
        closeSync(fd);
      }
    }
  } catch (error) {
    console.error(error);
    return;
  }
}

const params = {
  owner: "dhis2",
  repo: "dhis2-core",
  workflow_id: 10954,
};
fetchLatestRuns(params);
