import { Octokit } from "octokit";
import { openSync, closeSync, existsSync, writeFileSync} from "fs";

const octokit = new Octokit({
  auth: process.env.GITHUB_TOKEN,
});

async function fetchAllRuns(workflow_id : number) {
    octokit.hook.after("request", async (response, options) => {
        const ratelimit = `${response.headers["x-ratelimit-used"]}/${response.headers["x-ratelimit-limit"]}`
        console.log(`requested ${options.method} ${options.url}: ${response.status} ratelimit ${ratelimit}`)
    })
    octokit.hook.error("request", async (error, _) => {
        console.error(error)
    })

    try {
        const iterator = octokit.paginate.iterator(octokit.rest.actions.listWorkflowRuns, {
            owner: "dhis2",
            repo: "dhis2-core",
            event: "pull_request",
            status: "completed",
            workflow_id,
            per_page: 100,
        })

        for await (const { data: runs } of iterator) {
            for (const run of runs) {
                console.log("Run #%d", run.id);
                const file = `data/workflows/${workflow_id}/runs/${run.id}.json`
                if (existsSync(file)) {
                    console.log("run #%d exist already in %s", run.id, file)
                    continue
                }
                const fd = openSync(file, "w")
                writeFileSync(fd, JSON.stringify(run))
                closeSync(fd)
            }
        }
    } catch (error) {
        console.error(error)
        return 
    }
}

const workflowId = 10954
fetchAllRuns(workflowId)
