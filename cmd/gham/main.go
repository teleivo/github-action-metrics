// Command gham fetches GitHub Actions workflow data and indexes it in Elasticsearch.
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/teleivo/github-action-metrics/internal/cli"
)

var version = "dev"

// errFlagParse is a sentinel error indicating flag parsing failed.
// The flag package already printed the error, so main should not print again.
var errFlagParse = errors.New("flag parse error")

func main() {
	code, err := run(os.Args, os.Stdout, os.Stderr)
	if err != nil && err != errFlagParse {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	os.Exit(code)
}

func run(args []string, w io.Writer, wErr io.Writer) (int, error) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if len(args) < 2 {
		usage(wErr)
		return 2, nil
	}

	switch args[1] {
	case "fetch":
		return cli.HandleFetch(ctx, args[2:], wErr)
	case "index":
		return cli.HandleIndex(ctx, args[2:], wErr)
	case "version":
		fmt.Fprintln(w, version)
		return 0, nil
	case "-h", "--help", "help":
		usage(wErr)
		return 0, nil
	default:
		return 2, fmt.Errorf("unknown command: %s", args[1])
	}
}

func usage(w io.Writer) {
	fmt.Fprintln(w, `gham - GitHub Action Metrics

Fetch GitHub Actions workflow data and index it in Elasticsearch for analysis.

Usage:
  gham <command> [options]

Commands:
  fetch     Fetch workflow runs and jobs from GitHub
  index     Index stored data in Elasticsearch
  version   Print version information

Run 'gham <command> -h' for more information on a command.`)
}
