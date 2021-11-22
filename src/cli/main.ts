import { Command } from "commander";

import { makeFetchCommand } from "./fetch";

export function cli(args: string[]) {
  const program = new Command();
  program.version("0.0.1").showHelpAfterError();
  program.addCommand(makeFetchCommand());
  program.parse(args);
}
