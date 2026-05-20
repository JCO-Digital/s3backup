# S3 Backup Utility

Version: 1.2.0

A Go-based CLI utility for managing and downloading backups from S3-compatible storage (specifically optimized for Hetzner Object Storage). This tool automatically identifies the latest backups for various sites, manages local disk space, and provides both interactive and non-interactive download modes.

## Features

- **Automatic Grouping**: Groups S3 objects into "sites" based on the prefix (e.g., `site-name/backup.tar.gz`).
- **Latest Backup Detection**: Automatically finds the most recent database (`.sql.gz`) and file backups for each site using S3 metadata.
- **Space Management**: Monitors local disk space and automatically prunes oldest local backups to ensure a configurable amount of free space.
- **Interactive Mode**: Select sites to download using a terminal-based multi-select menu.
- **Cron Friendly**: A non-interactive mode (`--all`) perfect for scheduled backup tasks.
- **Status Overview**: Quickly view local storage statistics and age of backups.

## Installation

Ensure you have Go 1.21 or later installed.

```bash
git clone https://github.com/JCO-Digital/s3backup
cd s3backup
make build
```

The binary will be available in `./bin/s3backup`.

## Configuration

The utility looks for a configuration file at `~/.config/s3backup/config.toml` or `./config.toml`.

### Example `config.toml`

```toml
endpoint = "https://nbg1.your-objectstorage.com"
region = "us-east-1"
access_key = "YOUR_ACCESS_KEY"
secret_key = "YOUR_SECRET_KEY"
bucket = "my-backups"
target_path = "./backups"
max_age = 7        # Only consider backups from the last 7 days
keep_free = "100G" # Keep at least 100GB free on the target drive (supports "10%")
```

## Usage

### List Backups

List all detected sites in the S3 bucket and their latest backup timestamps:

```bash
./bin/s3backup list
```

### Download Backups

Download the latest backups for one or more specific sites:

```bash
./bin/s3backup download site-name
./bin/s3backup download site-one site-two
```

Download all latest backups (non-interactive):

```bash
./bin/s3backup download --all
```

Run interactively to select sites:

```bash
./bin/s3backup download
```

### Check Local Status

View information about local backups and disk space:

```bash
./bin/s3backup status
```

## Space Management Logic

When a download is initiated, the utility:

1. Calculates the total size of the files to be downloaded.
2. Checks available space on the `target_path` drive.
3. If `Available - Download Size < keep_free`, it deletes the oldest files in the `target_path` until the requirement is met.
4. If it cannot free up enough space, it exits with an error before starting the download.

## License

GPL-3.0
