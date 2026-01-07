import {
  GetObjectCommand,
  NoSuchKey,
  S3Client,
  S3ServiceException,
  paginateListObjectsV2,
} from "@aws-sdk/client-s3";
import { config } from "./main";
import { s3Object, objectSchema } from "./types";
import { basename, dirname, join } from "path";
import { createWriteStream, existsSync, mkdirSync } from "fs";

function getClientConfig() {
  return {
    endpoint: config.endpoint,
    region: config.region,
    forcePathStyle: true,
    credentials: {
      accessKeyId: config.accessKey,
      secretAccessKey: config.secretKey,
    },
  };
}

export async function listObjects(
  bucketName: string,
  pageSize: number = 100,
): Promise<s3Object[]> {
  const client = new S3Client(getClientConfig());
  const objects: s3Object[] = [];
  try {
    const paginator = paginateListObjectsV2(
      { client, pageSize },
      { Bucket: bucketName },
    );

    for await (const page of paginator) {
      if (page.Contents) {
        for (const item of page.Contents) {
          if (item.Size && item.Size > 0) {
            objects.push(objectSchema.parse(item));
          }
        }
      }
    }
  } catch (caught) {
    if (
      caught instanceof S3ServiceException &&
      caught.name === "NoSuchBucket"
    ) {
      console.error(
        `Error from S3 while listing objects for "${bucketName}". The bucket doesn't exist.`,
      );
    } else if (caught instanceof S3ServiceException) {
      console.error(
        `Error from S3 while listing objects for "${bucketName}".  ${caught.name}: ${caught.message}`,
      );
    } else {
      throw caught;
    }
  }
  return objects;
}

export async function downloadObject(bucketName: string, key: string) {
  const client = new S3Client(getClientConfig());

  if (!config.targetPath) {
    throw new Error("Target path is not defined");
  }

  try {
    const response = await client.send(
      new GetObjectCommand({
        Bucket: bucketName,
        Key: key,
      }),
    );

    if (!response.Body) {
      throw new Error("Response body is undefined");
    }

    const target = join(config.targetPath, key);

    if (existsSync(target)) {
      console.error(`File ${key} already exists in ${target}`);
      return;
    }

    // Create containing folder if it doesn't exist.
    const folder = dirname(target);
    if (!existsSync(folder)) {
      console.error(`Creating folder ${folder}`);
      mkdirSync(folder, { recursive: true });
    }

    console.log(`Downloading ${key} to ${target}`);

    // Pipe the response body directly to a file stream
    const fileStream = createWriteStream(target);
    response.Body.pipe(fileStream);

    return new Promise((resolve, reject) => {
      fileStream.on("finish", resolve);
      fileStream.on("error", reject);
    });
  } catch (caught) {
    if (caught instanceof NoSuchKey) {
      console.error(
        `Error from S3 while getting object "${key}" from "${bucketName}". No such key exists.`,
      );
    } else if (caught instanceof S3ServiceException) {
      console.error(
        `Error from S3 while getting object from ${bucketName}.  ${caught.name}: ${caught.message}`,
      );
    } else {
      throw caught;
    }
  }
}
