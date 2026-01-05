import {
  S3Client,
  S3ServiceException,
  paginateListObjectsV2,
} from "@aws-sdk/client-s3";
import { config } from "./main";

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

export async function listObjects(bucketName: string, pageSize: number = 100) {
  const client = new S3Client(getClientConfig());
  /** @type {string[][]} */
  const objects = [];
  try {
    const paginator = paginateListObjectsV2(
      { client, pageSize },
      { Bucket: bucketName },
    );

    for await (const page of paginator) {
      for (const item of page.Contents) {
        console.log(item);
      }
      //objects.push(page.Contents.map((o) => o.Key));
    }
    //console.log(objects);
    /*objects.forEach((objectList, pageNum) => {
    	console.log(
        `Page ${pageNum + 1}\n------\n${objectList.map((o) => `• ${o}`).join("\n")}\n`,
      );
      });*/
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
}
