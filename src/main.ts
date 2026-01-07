#!/usr/bin/env node
import { readConfigFile } from "./config";
import p from "../package.json";
import { downloadObjects, getLatestObjects } from "./process";

export const config = readConfigFile();

async function main() {
  console.error(`Version: ${p.version}`);

  const latest = await getLatestObjects();

  await downloadObjects(latest);
}

main();
