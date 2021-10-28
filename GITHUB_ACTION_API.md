# GitHub Action API

This is a collection of requests I found helpful ☺️

## How many runs where there on a specific day?

curl https://api.github.com/repos/dhis2/dhis2-core/actions/workflows/10954/runs\?event\=pull_request\&status\=completed\&created\=2021-10-26 | jq .total_count
