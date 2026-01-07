import { z } from "zod";

export const runtimeSchema = z.object({
  configDir: z.string(),
  cacheDir: z.string(),
  dataDir: z.string(),
  nodePath: z.string().default(""),
  scriptPath: z.string().default(""),
});

export type s3Runtime = z.infer<typeof runtimeSchema>;

export const configSchema = z.object({
  targetPath: z.string().default(""),
  bucket: z.string().default(""),
  endpoint: z.string().default(""),
  region: z.string().default(""),
  accessKey: z.string().default(""),
  secretKey: z.string().default(""),
});

export type s3Config = z.infer<typeof configSchema>;

export const objectSchema = z.object({
  Key: z.string(),
  LastModified: z.date(),
  ETag: z.string(),
  Size: z.number(),
  StorageClass: z.string(),
});

export type s3Object = z.infer<typeof objectSchema>;

export const latestSchema = z.object({
  site: z.string(),
  db: objectSchema.nullable(),
  files: objectSchema.nullable(),
});

export type s3Latest = z.infer<typeof latestSchema>;
