import { Octokit } from "octokit";
import { openSync, closeSync, existsSync, writeFileSync} from "fs";

const octokit = new Octokit({
  auth: process.env.GITHUB_TOKEN,
});

async function fetchRuns() {
    const file = "data/runs.json"
    if (existsSync(file)) {
        console.log(`runs exist already in ${file}`)
        return
    }

    let response
    try {
        response = await octokit.request('GET /repos/{owner}/{repo}/actions/runs', {
            owner: 'dhis2',
            repo: 'dhis2-core',
            event: 'pull_request',
            status: 'completed'
        })
    } catch (error) {
        console.error(error)
        return 
    }

    if (response.status != 200){
        console.log(`failed to fetch runs API responded with status ${response.status}`)
        return
    }

    const fd = openSync(file, "w")
    writeFileSync(fd, JSON.stringify(response.data))
    closeSync(fd)
}

fetchRuns()
