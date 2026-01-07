# s3backup

`s3backup` is a command-line tool designed to fetch and store the latest backup files from an Amazon S3 bucket locally. It is specifically tailored to work with backups created by the [SpinupWP](https://spinupwp.com/) backup system, providing an extra layer of redundancy by keeping a local copy of your most recent backups.

## Features

- Downloads the latest backup file from a specified S3 bucket
- Stores backups locally for redundancy
- Simple configuration via a TOML file
- Designed to be run via cron or other schedulers
- Compatible with custom S3 endpoints (e.g., Hetzner, MinIO, Wasabi)

## Installation

Clone this repository and build the tool according to your platform’s requirements.

```sh
git clone https://github.com/yourusername/s3backup.git
cd s3backup
make
```

Configuration

Create a configuration file at `~/.config/s3backup/config.toml` with the following structure:

```toml
targetPath = "/path/to/your/local/directory"
bucket     = "your-s3-bucket-name"
endpoint   = "https://s3.amazonaws.com" # or your custom endpoint
region     = "us-east-1"
accessKey  = "YOUR_AWS_ACCESS_KEY"
secretKey  = "YOUR_AWS_SECRET_KEY"
```

## Usage

Run the tool from the command line:

```sh
s3backup
```

This will download the latest backup file from your configured S3 bucket and store it locally.

## Scheduling

`s3backup` does not handle scheduling internally. To automate backups, add it to your `cron` jobs or use another scheduler:

```cron
0 2 * * * /usr/local/bin/s3backup
```
