import * as fs from "fs";
import * as path from "path";

import { Octokit } from "octokit";
import { RestEndpointMethodTypes } from "@octokit/plugin-rest-endpoint-methods";

async function fetchJobs(
  params: RestEndpointMethodTypes["actions"]["listJobsForWorkflowRun"]["parameters"],
  workflowId: number,
  directory: string,
  token: string
) {
  const octokit = new Octokit({
    auth: token,
  });

  const file = path.join(
    directory,
    `/workflows/${workflowId}/jobs/${params.run_id}.json`
  );
  if (fs.existsSync(file)) {
    console.log("jobs of run #%d already exist in %s", params.run_id, file);
    return;
  }

  let response;
  try {
    // https://docs.github.com/en/rest/reference/actions#list-jobs-for-a-workflow-run
    // currently I just take the 'latest' but there could be multiple
    // attemps see filter
    response = await octokit.rest.actions.listJobsForWorkflowRun(params);
    console.log("Fetched jobs for run #%d", params.run_id);
  } catch (error) {
    console.error(error);
    return;
  }
  if (response.status != 200) {
    console.log(
      `failed to fetch jobs API responded with status ${response.status}`
    );
    return;
  }

  const fd = fs.openSync(file, "w");
  fs.writeFileSync(fd, JSON.stringify(response.data));
  // TODO ensure I always close the file
  fs.closeSync(fd);
}

async function fetch(
  repo: string,
  owner: string,
  workflowId: number,
  directory: string,
  token: string
) {
  // TODO handle that directory not being present and give meaningful
  // error message telling the user to fetch runs before
  const file = path.join(directory, `/workflows/${workflowId}/runs/`);
  const dir = fs.opendirSync(file);
  for await (const dirent of dir) {
    if (!dirent.isFile()) {
      continue;
    }
    const runId = Number(path.parse(dirent.name).name);
    if (Number.isNaN(runId)) {
      console.log(`failed to parse runId from file ${dirent.name}`);
      continue;
    }
    fetchJobs({ owner, repo, run_id: runId }, workflowId, directory, token);
  }
}

function fetchAllJobs(
  repo: string,
  owner: string,
  workflowId: number,
  directory: string,
  token: string
) {
  fetch(repo, owner, workflowId, directory, token);
}
export default fetchAllJobs;
