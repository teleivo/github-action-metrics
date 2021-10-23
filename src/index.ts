import { Octokit } from "octokit";
import { openSync, closeSync, existsSync, writeFileSync} from "fs";

const octokit = new Octokit({
  auth: process.env.GITHUB_TOKEN,
});

async function fetchWorkflows() {
    const file = "data/workflows.json"
    if (existsSync(file)) {
        console.log(`workflows exist already in ${file}`)
        return
    }

    let response
    try {
        response = await octokit.request('GET /repos/{owner}/{repo}/actions/workflows', {
            owner: 'dhis2',
            repo: 'dhis2-core'
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

fetchWorkflows()
