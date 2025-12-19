// Command gham fetches GitHub Actions workflow data and indexes it in Elasticsearch.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/teleivo/github-action-metrics/internal/cli"
)

var version = "dev"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "fetch":
		cli.HandleFetch(ctx, os.Args[2:])
	case "index":
		cli.HandleIndex(ctx, os.Args[2:])
	case "version":
		fmt.Println(version)
	case "-h", "--help", "help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`gham - GitHub Action Metrics

Fetch GitHub Actions workflow data and index it in Elasticsearch for analysis.

Usage:
  gham <command> [options]

Commands:
  fetch     Fetch workflow runs and jobs from GitHub
  index     Index stored data in Elasticsearch
  version   Print version information

Run 'gham <command> -h' for more information on a command.`)
}
