import fs from "fs";
import path from "path";

import { Client } from "@elastic/elasticsearch";

async function* generatorRuns(workflowId: number, srcDir: string) {
  const file = path.join(srcDir, `/workflows/${workflowId}/runs/`);
  const dir = fs.opendirSync(file);
  for await (const dirent of dir) {
    if (!dirent.isFile()) {
      continue;
    }

    let result;
    try {
      const run = JSON.parse(
        fs.readFileSync(path.join(file, dirent.name), "utf8")
      );
      const jobs = JSON.parse(
        fs.readFileSync(
          path.join(srcDir, `/workflows/${workflowId}/jobs/${run.id}.json`),
          "utf8"
        )
      );
      result = { ...run, ...runDuration(jobs) };
    } catch (err) {
      console.error(err);
      continue;
    }

    yield result;
  }
}

async function* generatorJobs(workflowId: number, srcDir: string) {
  const file = path.join(srcDir, `/workflows/${workflowId}/jobs/`);
  const dir = fs.opendirSync(file);
  for await (const dirent of dir) {
    if (!dirent.isFile()) {
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
          job_url: job.url,
          job_html_url: job.html_url,
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

export async function indexRuns(
  node: string,
  username: string,
  password: string,
  workflowId: number,
  srcDir: string
) {
  const client = new Client({ node, auth: { username, password } });
  const result = await client.helpers.bulk({
    datasource: generatorRuns(workflowId, srcDir),
    // TODO add a type for the doc?
    onDocument(doc: any) {
      return {
        index: { _index: "runs", _id: doc.id },
      };
    },
  });

  console.log("indexed runs", result);
}

export async function indexJobs(
  node: string,
  username: string,
  password: string,
  workflowId: number,
  srcDir: string
) {
  const client = new Client({ node, auth: { username, password } });
  const result = await client.helpers.bulk({
    datasource: generatorJobs(workflowId, srcDir),
    // TODO add a type for the doc?
    onDocument(doc: any) {
      return {
        index: { _index: "jobs", _id: doc.id },
      };
    },
  });

  console.log("indexed jobs", result);
}

export async function indexSteps(
  node: string,
  username: string,
  password: string,
  workflowId: number,
  srcDir: string
) {
  const client = new Client({ node, auth: { username, password } });
  const result = await client.helpers.bulk({
    datasource: generatorSteps(workflowId, srcDir),
    // TODO add a type for the doc?
    onDocument(doc: any) {
      return {
        index: { _index: "steps", _id: doc.job_id + "-" + doc.number },
      };
    },
  });

  console.log("indexed steps", result);
}

type RunDuration = {
  jobs_started_at: string;
  jobs_started_at_id: string;
  jobs_started_at_name: string;
  jobs_started_at_url: string;
  jobs_started_at_html_url: string;
  jobs_completed_at: string;
  jobs_completed_at_id: string;
  jobs_completed_at_name: string;
  jobs_completed_at_url: string;
  jobs_completed_at_html_url: string;
};

export function runDuration(jobs: any): RunDuration | {} {
  return jobs.jobs.reduce((run: any, j: any) => {
    if (
      !run.jobs_started_at ||
      new Date(j.started_at) < new Date(run.jobs_started_at)
    ) {
      run.jobs_started_at = j.started_at;
      run.jobs_started_at_id = j.id;
      run.jobs_started_at_name = j.name;
      run.jobs_started_at_url = j.url;
      run.jobs_started_at_html_url = j.html_url;
    }
    if (
      !run.jobs_completed_at ||
      new Date(j.completed_at) > new Date(run.jobs_completed_at)
    ) {
      run.jobs_completed_at = j.completed_at;
      run.jobs_completed_at_id = j.id;
      run.jobs_completed_at_name = j.name;
      run.jobs_completed_at_url = j.url;
      run.jobs_completed_at_html_url = j.html_url;
    }
    return run;
  }, {});
}
