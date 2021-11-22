import * as path from "path";
import * as fs from "fs";

import { Command } from "commander";

import { indexSteps, indexJobs, indexRuns } from "../elastic";

export function makeIndexCommand(): Command {
  const index = new Command("index").description(
    "Index stored GitHub workflow runs and jobs in Elasticsearch"
  );
  index
    .command("runs")
    .description("Index GitHub workflow runs in Elasticsearch")
    .requiredOption(
      "-u, --url <value>",
      "Elasticsearch URL where data will be indexed"
    )
    .requiredOption(
      "-w, --workflow-id <value>",
      "Workflow id of GitHub action",
      parseInt
    )
    .requiredOption(
      "-s, --source <value>",
      "Directory where GitHub action payloads are stored"
    )
    .requiredOption(
      "-u, --user <value>",
      "Elasticsearch basic authentication user"
    )
    .requiredOption(
      "-p, --password <value>",
      "Elasticsearch basic authentication password"
    )
    .action(executeRuns);

  index
    .command("jobs")
    .description("Index GitHub workflow jobs in Elasticsearch")
    .requiredOption(
      "-u, --url <value>",
      "Elasticsearch URL where data will be indexed"
    )
    .requiredOption(
      "-w, --workflow-id <value>",
      "Workflow id of GitHub action",
      parseInt
    )
    .requiredOption(
      "-s, --source <value>",
      "Directory where GitHub action payloads are stored"
    )
    .requiredOption(
      "-u, --user <value>",
      "Elasticsearch basic authentication user"
    )
    .requiredOption(
      "-p, --password <value>",
      "Elasticsearch basic authentication password"
    )
    .action(executeJobs);

  index
    .command("steps")
    .description("Index GitHub workflow steps in Elasticsearch")
    .requiredOption(
      "-u, --url <value>",
      "Elasticsearch URL where data will be indexed"
    )
    .requiredOption(
      "-w, --workflow-id <value>",
      "Workflow id of GitHub action",
      parseInt
    )
    .requiredOption(
      "-s, --source <value>",
      "Directory where GitHub action payloads are stored"
    )
    .requiredOption(
      "-u, --user <value>",
      "Elasticsearch basic authentication user"
    )
    .requiredOption(
      "-p, --password <value>",
      "Elasticsearch basic authentication password"
    )
    .action(executeSteps);

  index
    .command("all")
    .description("Index GitHub workflow runs, jobs and steps in Elasticsearch")
    .requiredOption(
      "-u, --url <value>",
      "Elasticsearch URL where data will be indexed"
    )
    .requiredOption(
      "-w, --workflow-id <value>",
      "Workflow id of GitHub action",
      parseInt
    )
    .requiredOption(
      "-s, --source <value>",
      "Directory where GitHub action payloads are stored"
    )
    .requiredOption(
      "-u, --user <value>",
      "Elasticsearch basic authentication user"
    )
    .requiredOption(
      "-p, --password <value>",
      "Elasticsearch basic authentication password"
    )
    .action(executeAll);

  return index;
}

async function executeRuns(options: any): Promise<void> {
  let dir: string;
  try {
    dir = path.resolve(options.source);
    if (!fs.lstatSync(dir).isDirectory()) {
      console.error(`${options.source} must be a directory`);
      process.exit(1);
    }
  } catch (err) {
    console.error(err);
    process.exit(1);
  }

  indexRuns(
    options.url,
    options.user,
    options.password,
    options.workflowId,
    dir
  );
}

async function executeJobs(options: any): Promise<void> {
  let dir: string;
  try {
    dir = path.resolve(options.source);
    if (!fs.lstatSync(dir).isDirectory()) {
      console.error(`${options.source} must be a directory`);
      process.exit(1);
    }
  } catch (err) {
    console.error(err);
    process.exit(1);
  }

  indexJobs(
    options.url,
    options.user,
    options.password,
    options.workflowId,
    dir
  );
}

async function executeSteps(options: any): Promise<void> {
  let dir: string;
  try {
    dir = path.resolve(options.source);
    if (!fs.lstatSync(dir).isDirectory()) {
      console.error(`${options.source} must be a directory`);
      process.exit(1);
    }
  } catch (err) {
    console.error(err);
    process.exit(1);
  }

  indexSteps(
    options.url,
    options.user,
    options.password,
    options.workflowId,
    dir
  );
}

async function executeAll(options: any): Promise<void> {
  let dir: string;
  try {
    dir = path.resolve(options.source);
    if (!fs.lstatSync(dir).isDirectory()) {
      console.error(`${options.source} must be a directory`);
      process.exit(1);
    }
  } catch (err) {
    console.error(err);
    process.exit(1);
  }

  console.log("Index runs");
  indexRuns(
    options.url,
    options.user,
    options.password,
    options.workflowId,
    dir
  );
  console.log("Index jobs");
  indexJobs(
    options.url,
    options.user,
    options.password,
    options.workflowId,
    dir
  );
  console.log("Index steps");
  indexSteps(
    options.url,
    options.user,
    options.password,
    options.workflowId,
    dir
  );
}
