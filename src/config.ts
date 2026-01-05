import { xdgConfig, xdgCache, xdgData } from "xdg-basedir";
import { join } from "path";
import { readFileSync, writeFileSync, existsSync, mkdirSync } from "fs";
import { parse as tomlParse } from "smol-toml";
import { configSchema, s3Config, runtimeSchema } from "./types";
import p from "../package.json";

export const runtimeData = runtimeSchema.parse({
  configDir: join(xdgConfig ?? "", p.name),
  cacheDir: join(xdgCache ?? "", p.name),
  dataDir: join(xdgData ?? "", p.name),
});

export function getConfigFilePath(): string {
  [runtimeData.configDir, runtimeData.cacheDir, runtimeData.dataDir].forEach(
    (path) => {
      if (!existsSync(path)) {
        console.debug(`Creating directory at ${path}`);
        mkdirSync(path, { recursive: true });
      }
    },
  );
  return join(runtimeData.configDir, "config.toml");
}

export function readConfigFile(): s3Config {
  const configPath = getConfigFilePath();
  if (!existsSync(configPath)) {
    console.debug(`Creating default config file at ${configPath}`);
    writeFileSync(configPath, "");
  }
  const configRaw = readFileSync(configPath, "utf8");
  const config = configSchema.parse(tomlParse(configRaw));

  return config;
}
