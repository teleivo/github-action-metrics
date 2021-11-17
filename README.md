# GitHub Action Metrics

[![main](https://github.com/teleivo/github-action-metrics/actions/workflows/main.yml/badge.svg)](https://github.com/teleivo/github-action-metrics/actions/workflows/main.yml)

Analyze your [GitHub Actions](https://github.com/features/actions)!

Do you wonder
- why PRs take so long to be checked by your GitHub actions?
- how often PRs fail to pass a certain workflow, job or step?
- if caching dependencies did make a job run faster?

These are just a fraction of questions you can find answers to using this
project :blush:

## Architecture

This project provides a CLI which is a small wrapper around [GitHub's
Octokit](https://github.com/octokit/octokit.js) library. It will fetch GitHub
action data from https://docs.github.com/en/rest/reference/actions and store
it in a place of your choice. You can for example store it in a Git repository
on GitHub itself.

You can then index the data into
[Elasticsearch](https://www.elastic.co/elasticsearch/). Create visualizations
(graphs, metrics, dashboards, ...) to answer a lot of questions you have about
your use of GitHub actions.

## Example

I started this project to analyze the test workflow we use at
[DHIS2](https://dhis2.org/about/). I wanted to know where time was spent, why
some test runs took 15min while others took 23min to finish. How could we
get faster feedback on PRs and reduce this variation in test duration?

Using this project I

- fetched GitHub action data every day using a scheduled GitHub action and
  stored it in a [GitHub
  repository](https://github.com/teleivo/dhis2-github-action-metrics)
- indexed and analyzed the data using Elasticsearch and Kibana
