"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const octokit_1 = require("octokit");
const fs_1 = require("fs");
const octokit = new octokit_1.Octokit();
async function fetchWorkflows() {
    const file = "data/workflows.json";
    const fd = (0, fs_1.openSync)(file, "w");
    const response = await octokit.request('GET /repos/{owner}/{repo}/actions/workflows', {
        owner: 'dhis2',
        repo: 'dhis2-core'
    });
    if (response.status === 200) {
        (0, fs_1.writeFileSync)(fd, JSON.stringify(response.data));
    }
}
fetchWorkflows();
