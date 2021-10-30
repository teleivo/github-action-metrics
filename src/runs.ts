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
export async function fetchRuns(
  repo: string,
  owner: string,
  workflowId: number,
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
    const params: RestEndpointMethodTypes["actions"]["listWorkflowRuns"]["parameters"] =
      {
        repo,
        owner,
        workflow_id: workflowId,
        event: "pull_request",
        status: "completed",
        per_page: 100,
      };
    if (created) {
      params.created = created;
    }

    const iterator = octokit.paginate.iterator(
      octokit.rest.actions.listWorkflowRuns,
      {
        ...params,
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

        let fd;
        try {
          fd = fs.openSync(file, "w");
          fs.writeFileSync(fd, JSON.stringify(run), "utf8");
        } finally {
          if (fd !== undefined) {
            fs.closeSync(fd);
          }
        }
      }
    }
  } catch (error) {
    console.error(error);
    return;
  }
}
