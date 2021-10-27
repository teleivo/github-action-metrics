import { Octokit } from "octokit";
import {
  openSync,
  opendirSync,
  closeSync,
  existsSync,
  writeFileSync,
} from "fs";
import * as path from "path";
import { RestEndpointMethodTypes } from "@octokit/plugin-rest-endpoint-methods";

const octokit = new Octokit({
  auth: process.env.GITHUB_TOKEN,
});

async function fetchJobs(
  workflowId: number,
  params: RestEndpointMethodTypes["actions"]["listJobsForWorkflowRun"]["parameters"]
) {
  const file = `data/workflows/${workflowId}/jobs/${params.run_id}.json`;
  if (existsSync(file)) {
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

  const fd = openSync(file, "w");
  writeFileSync(fd, JSON.stringify(response.data));
  closeSync(fd);
}

async function fetchAllJobs(owner: string, repo: string, workflowId: number) {
  const file = `data/workflows/${workflowId}/runs/`;
  const dir = opendirSync(file);
  for await (const dirent of dir) {
    if (!dirent.isFile()) {
      continue;
    }
    const runId = Number(path.parse(dirent.name).name);
    if (Number.isNaN(runId)) {
      console.log(`failed to parse runId from file ${dirent.name}`);
      continue;
    }
    fetchJobs(workflowId, { owner, repo, run_id: runId });
  }
}

const owner = "dhis2";
const repo = "dhis2-core";
const workflowId = 10954;
fetchAllJobs(owner, repo, workflowId);
