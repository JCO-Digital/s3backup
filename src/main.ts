#!/usr/bin/env node
import { readConfigFile } from "./config";
import p from "../package.json";
import { listObjects } from "./s3utils";

export const config = readConfigFile();

async function main() {
  console.error(`Version: ${p.version}`);
  listObjects("spinup");
}

main();
