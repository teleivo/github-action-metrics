import { Command } from "commander";

import { makeFetchCommand } from "./fetch";
import { makeIndexCommand } from "./index";

export function cli(args: string[]) {
  const program = new Command();
  program.version("0.0.1").showHelpAfterError();
  program.addCommand(makeFetchCommand());
  program.addCommand(makeIndexCommand());
  program.parse(args);
}
