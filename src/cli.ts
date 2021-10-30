import * as path from "path";
import * as fs from "fs";

import { Command } from "commander";

import { fetchRuns } from "./runs";
import fetchAllJobs from "./jobs";

export function cli(args: string[]) {
  const program = new Command();
  program.version("0.0.1").showHelpAfterError();

  program
    .command("runs")
    .description(
      "Fetch latest GitHub action runs for given workflow via https://docs.github.com/en/rest/reference/actions#list-workflow-runs"
    )
    .requiredOption("-r, --repo <value>", "GitHub repository")
    .requiredOption("-o, --owner <value>", "Owner of GitHub repository")
    .requiredOption(
      "-w, --workflow-id <value>",
      "Workflow id of GitHub action",
      parseInt
    )
    .requiredOption(
      "-d, --directory <value>",
      "Directory where GitHub action payloads will be stored"
    )
    .option(
      "-c, --created <value>",
      "Date the run was created in format like '2021-10-12' or '2021-10-29T22:40:19Z'"
    )
    .option("-t, --token <value>", "GitHub access token")
    .action(executeRuns);

  program
    .command("jobs")
    .description(
      "Fetch all GitHub action jobs of stored runs via https://docs.github.com/en/rest/reference/actions#get-a-job-for-a-workflow-run"
    )
    .requiredOption("-r, --repo <value>", "GitHub repository")
    .requiredOption("-o, --owner <value>", "Owner of GitHub repository")
    .requiredOption(
      "-w, --workflow-id <value>",
      "Workflow id of GitHub action",
      parseInt
    )
    .requiredOption(
      "-d, --directory <value>",
      "Directory where GitHub action payloads will be stored"
    )
    .option("-t, --token <value>", "GitHub access token")
    .description("fetch latest runs of given GitHub workflow")
    .action(executeJobs);
  program.parse(args);
}

function executeRuns(options: any): void {
  let dir: string;
  try {
    dir = path.resolve(options.directory);
    if (!fs.lstatSync(dir).isDirectory()) {
      console.error(`${options.directory} must be a directory`);
      process.exit(1);
    }
  } catch (err) {
    console.error(err);
    process.exit(1);
  }

  fetchRuns(
    options.repo,
    options.owner,
    options.workflowId,
    dir,
    options.created,
    options.token
  );
}

function executeJobs(options: any): void {
  let dir: string;
  try {
    dir = path.resolve(options.directory);
    if (!fs.lstatSync(dir).isDirectory()) {
      console.error(`${options.directory} must be a directory`);
      process.exit(1);
    }
  } catch (err) {
    console.error(err);
    process.exit(1);
  }
  fetchAllJobs(
    options.repo,
    options.owner,
    options.workflowId,
    dir,
    options.token
  );
}
