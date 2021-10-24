import { Octokit } from "octokit";
import { openSync, closeSync, existsSync, writeFileSync} from "fs";

const octokit = new Octokit({
  auth: process.env.GITHUB_TOKEN,
});

async function fetchJobs(workflowId : number, runId : number) {
    const file = `data/workflows/${workflowId}/jobs/${runId}.json`
    if (existsSync(file)) {
        console.log("jobs of run #%d already exist in %s", runId, file)
        return
    }

    let response
    try {
        // https://docs.github.com/en/rest/reference/actions#list-jobs-for-a-workflow-run
        // currently I just take the 'latest' but there could be multiple
        // attemps see filter
        response = await octokit.rest.actions.listJobsForWorkflowRun({
            owner: "dhis2",
            repo: "dhis2-core",
            run_id: runId,
        })
    } catch (error) {
        console.error(error)
        return 
    }
    if (response.status != 200){
        console.log(`failed to fetch workflows API responded with status ${response.status}`)
        return
    }

    const fd = openSync(file, "w")
    writeFileSync(fd, JSON.stringify(response.data))
    closeSync(fd)
}

const workflowId = 10954
// TODO call fetchJobs for every ${runId}.json file in
// data/workflows/${workflowId}/...json
// this is just a sample
const runId = 1233841997
fetchJobs(workflowId, runId)
