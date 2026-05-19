package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/niklas/s3backup/internal/config"
	"github.com/niklas/s3backup/internal/s3"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all detected sites and their latest backups",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		cobra.CheckErr(err)

		client, err := s3.NewClient(cfg)
		cobra.CheckErr(err)

		sites, err := client.GetLatestBackups(context.Background())
		cobra.CheckErr(err)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "SITE\tDB LATEST\tFILES LATEST\tTOTAL SIZE")
		var totalDB int64
		var totalFiles int64

		for _, site := range sites {
			dbDate := "-"
			if site.DB != nil {
				dbDate = site.DB.LastModified.Format("2006-01-02 15:04:05")
				totalDB += site.DB.Size
			}
			filesDate := "-"
			if site.Files != nil {
				filesDate = site.Files.LastModified.Format("2006-01-02 15:04:05")
				totalFiles += site.Files.Size
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				site.Name,
				dbDate,
				filesDate,
				formatSize(site.TotalSize()))
		}
		fmt.Fprintln(w, "\t\t\t") // spacer
		fmt.Fprintf(w, "TOTAL\t%s\t%s\t%s\n",
			formatSize(totalDB),
			formatSize(totalFiles),
			formatSize(totalDB+totalFiles))
		w.Flush()
	},
}

func formatSize(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func init() {
	rootCmd.AddCommand(listCmd)
}
