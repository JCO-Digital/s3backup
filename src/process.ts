import { downloadObject, listObjects } from "./s3utils";
import { config } from "./main";
import { s3Latest } from "./types";
import { maxAge } from "./constants";
import { displayFileSize } from "./utils";

export async function getLatestObjects(): Promise<s3Latest[]> {
  let latest: s3Latest[] = [];
  const objects = await listObjects(config.bucket);
  for (const object of objects) {
    // Skip objects older than maxAge
    if (Date.now() - object.LastModified.getTime() > maxAge) continue;

    const [site, file] = object.Key.split("/");
    let latestItem = null;
    for (const l of latest) {
      if (l.site === site) {
        latestItem = l;
        break;
      }
    }
    if (latestItem === null) {
      latestItem = {
        site,
        db: null,
        files: null,
      };
      latest.push(latestItem);
    }
    const isDb = file.endsWith(".sql.gz");
    if (isDb) {
      if (
        latestItem.db === null ||
        latestItem.db.LastModified < object.LastModified
      ) {
        latestItem.db = object;
      }
    } else {
      if (
        latestItem.files === null ||
        latestItem.files.LastModified < object.LastModified
      ) {
        latestItem.files = object;
      }
    }
  }
  return latest;
}

export async function downloadObjects(latest: s3Latest[]) {
  let dbSize = 0;
  let fileSize = 0;

  for (const item of latest) {
    if (item.db) {
      dbSize += item.db.Size;
    }
    if (item.files) {
      fileSize += item.files.Size;
    }
  }
  console.log(`Total database size: ${displayFileSize(dbSize)}`);
  console.log(`Total file size: ${displayFileSize(fileSize)}`);

  /*
	for (const item of latest) {
		if (item.db) {
			await downloadObject(config.bucket, item.db.Key);
		}
		if (item.files) {
			await downloadObject(config.bucket, item.files.Key);
		}
	}
	*/
}
