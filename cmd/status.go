package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/niklas/s3backup/internal/config"
	"github.com/niklas/s3backup/internal/utils"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status of local backups",
	Long:  `Displays information about locally stored backups, including total size, file count, and age range.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		cobra.CheckErr(err)

		if cfg.TargetPath == "" {
			cobra.CheckErr(fmt.Errorf("target_path must be set in config"))
		}

		absPath, err := filepath.Abs(cfg.TargetPath)
		cobra.CheckErr(err)

		var totalSize int64
		var fileCount int
		var oldest time.Time
		var newest time.Time
		var first = true

		err = filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Skip errors
			}
			if !info.IsDir() {
				totalSize += info.Size()
				fileCount++

				modTime := info.ModTime()
				if first {
					oldest = modTime
					newest = modTime
					first = false
				} else {
					if modTime.Before(oldest) {
						oldest = modTime
					}
					if modTime.After(newest) {
						newest = modTime
					}
				}
			}
			return nil
		})
		cobra.CheckErr(err)

		disk, err := utils.GetDiskSpace(absPath)
		if err != nil {
			// Try parent if target doesn't exist yet
			disk, _ = utils.GetDiskSpace(filepath.Dir(absPath))
		}

		fmt.Printf("Local Backup Status (%s):\n", absPath)
		fmt.Printf("  Files:            %d\n", fileCount)
		fmt.Printf("  Total Size:       %s\n", formatSize(totalSize))

		if fileCount > 0 {
			fmt.Printf("  Oldest Backup:    %s (%s ago)\n", oldest.Format("2006-01-02 15:04:05"), formatDuration(time.Since(oldest)))
			fmt.Printf("  Most Recent:      %s (%s ago)\n", newest.Format("2006-01-02 15:04:05"), formatDuration(time.Since(newest)))
		} else {
			fmt.Printf("  Backups:          None found\n")
		}

		if disk.All > 0 {
			minFree, _ := utils.ParseSize(cfg.KeepFree, disk.All)
			fmt.Printf("\nDisk Space:\n")
			fmt.Printf("  Total:            %s\n", formatSize(int64(disk.All)))
			fmt.Printf("  Available:        %s\n", formatSize(int64(disk.Free)))
			fmt.Printf("  Keep Free Limit:  %s (%s)\n", formatSize(int64(minFree)), cfg.KeepFree)
		}
	},
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh", days, hours)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
