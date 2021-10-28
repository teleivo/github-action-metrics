import * as path from "path";
import * as fs from "fs";

import { Command } from "commander";

import fetchLatestRun from "./runs";

export function cli(args: string[]) {
  const program = new Command();
  program.version("0.0.1").showHelpAfterError();

  // TODO how can I reuse the options? I want them on the root and accessible
  // in the subcommands. Using the command parameter in the action handler
  // did not work.
  program
    .command("runs")
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
    .action((options) => {
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
      fetchLatestRun(
        options.repo,
        options.owner,
        options.workflowId,
        dir,
        options.token
      );
    });
  program.parse(args);
}
