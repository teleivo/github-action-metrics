import * as path from "path";
import * as fs from "fs";

import { Command } from "commander";

import { fetchLatestRuns, fetchRunsCreatedOn } from "./runs";
import fetchAllJobs from "./jobs";

export function cli(args: string[]) {
  const program = new Command();
  program.version("0.0.1").showHelpAfterError();

  // TODO how can I reuse the options? I want them on the root and accessible
  // in the subcommands. Using the command parameter in the action handler
  // did not work.
  const runs = program
    .command("runs")
    .description("fetch latest GitHub action runs for given workflow");

  runs
    .command("latest", { isDefault: true })
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
    .action(executeRunsLatest);
  runs
    .command("created")
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
    .requiredOption(
      "-c, --created <value>",
      "Date the run was created in format like '2021-10-12'"
    )
    .option("-t, --token <value>", "GitHub access token")
    .action(executeRunsCreated);

  program
    .command("jobs")
    .description("fetch all GitHub action jobs of stored runs")
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

function executeRunsCreated(options: any): void {
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
  fetchRunsCreatedOn(
    options.repo,
    options.owner,
    options.workflowId,
    dir,
    options.created,
    options.token
  );
}

function executeRunsLatest(options: any): void {
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
  fetchLatestRuns(
    options.repo,
    options.owner,
    options.workflowId,
    dir,
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
