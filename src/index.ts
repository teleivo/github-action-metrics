import { Octokit } from "octokit";
import { openSync, writeFileSync} from "fs";

const octokit = new Octokit();

async function fetchWorkflows() {
    const file = "data/workflows.json"

    const fd=openSync(file, "w")

    const response = await octokit.request('GET /repos/{owner}/{repo}/actions/workflows', {
        owner: 'dhis2',
        repo: 'dhis2-core'
    })

    if (response.status === 200){
        writeFileSync(fd, JSON.stringify(response.data))
    }
}

fetchWorkflows()
