import fs from "fs";
import path from "path";

import { Client } from "@elastic/elasticsearch";

// TODO expose this in the CLI
// gm ingest runs
// gm ingest jobs

// rename others into subcommands?
// gm fetch runs
// gm fetch jobs

// TODO ingest jobs
// note jobs have a total_count field with potentially multiple jobs per run.
// Do I want to ingest it as is? I feel like I would rather want them to be
// ingested separately so the job will have an appropriate mapping in ES as
// well :) each job already has the run_id as a field so I can tie it back
// together with its run :)
async function* generatorRuns(workflowId: number, srcDir: string) {
  const file = path.join(srcDir, `/workflows/${workflowId}/runs/`);
  const dir = fs.opendirSync(file);
  for await (const dirent of dir) {
    if (!dirent.isFile()) {
      continue;
    }
    const runId = Number(path.parse(dirent.name).name);
    if (Number.isNaN(runId)) {
      console.log(`failed to parse runId from file ${dirent.name}`);
      continue;
    }
    const data = fs.readFileSync(path.join(file, dirent.name), "utf8");
    yield JSON.parse(data);
  }
}

async function* generatorJobs(workflowId: number, srcDir: string) {
  const file = path.join(srcDir, `/workflows/${workflowId}/jobs/`);
  const dir = fs.opendirSync(file);
  for await (const dirent of dir) {
    if (!dirent.isFile()) {
      continue;
    }
    const runId = Number(path.parse(dirent.name).name);
    if (Number.isNaN(runId)) {
      console.log(`failed to parse runId from file ${dirent.name}`);
      continue;
    }
    const data = fs.readFileSync(path.join(file, dirent.name), "utf8");
    const jobs = JSON.parse(data);
    for (const job of jobs.jobs) {
      yield job;
    }
  }
}
async function* generatorSteps(workflowId: number, srcDir: string) {
  const file = path.join(srcDir, `/workflows/${workflowId}/jobs/`);
  const dir = fs.opendirSync(file);
  for await (const dirent of dir) {
    if (!dirent.isFile()) {
      continue;
    }
    const runId = Number(path.parse(dirent.name).name);
    if (Number.isNaN(runId)) {
      console.log(`failed to parse runId from file ${dirent.name}`);
      continue;
    }
    const data = fs.readFileSync(path.join(file, dirent.name), "utf8");
    const jobs = JSON.parse(data);
    for (const job of jobs.jobs) {
      for (const step of job.steps) {
        yield {
          ...step,
          job_id: job.id,
          job_name: job.name,
          run_id: job.run_id,
          run_url: job.run_url,
          run_html_url: job.html_url,
          run_attempt: job.run_attempt,
          head_sha: job.head_sha,
        };
      }
    }
  }
}

async function ingestRuns(node: string, workflowId: number, srcDir: string) {
  const client = new Client({ node });
  const result = await client.helpers.bulk({
    datasource: generatorRuns(workflowId, srcDir),
    // TODO add a type for the doc?
    onDocument(doc: any) {
      return {
        index: { _index: "runs", _id: doc.id },
      };
    },
  });

  console.log(result);
}

async function ingestJobs(node: string, workflowId: number, srcDir: string) {
  const client = new Client({ node });
  const result = await client.helpers.bulk({
    datasource: generatorJobs(workflowId, srcDir),
    // TODO add a type for the doc?
    onDocument(doc: any) {
      return {
        index: { _index: "jobs", _id: doc.id },
      };
    },
  });

  console.log(result);
}

async function ingestSteps(node: string, workflowId: number, srcDir: string) {
  const client = new Client({ node });
  const result = await client.helpers.bulk({
    datasource: generatorSteps(workflowId, srcDir),
    // TODO add a type for the doc?
    onDocument(doc: any) {
      return {
        index: { _index: "steps", _id: doc.job_id + "-" + doc.number },
      };
    },
  });

  console.log(result);
}

ingestRuns("http://localhost:9200", 10954, "/home/ivo/code/dhis2/metrics/data");
ingestJobs("http://localhost:9200", 10954, "/home/ivo/code/dhis2/metrics/data");
ingestSteps(
  "http://localhost:9200",
  10954,
  "/home/ivo/code/dhis2/metrics/data"
);
