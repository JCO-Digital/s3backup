package cmd

import (
	"context"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/niklas/s3backup/internal/backup"
	"github.com/niklas/s3backup/internal/config"
	"github.com/niklas/s3backup/internal/s3"
	"github.com/spf13/cobra"
)

var downloadAll bool

var downloadCmd = &cobra.Command{
	Use:   "download [site]",
	Short: "Download latest backups for a site or all sites",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig()
		cobra.CheckErr(err)

		if cfg.TargetPath == "" {
			cobra.CheckErr(fmt.Errorf("target_path must be set in config"))
		}

		client, err := s3.NewClient(cfg)
		cobra.CheckErr(err)

		sites, err := client.GetLatestBackups(context.Background())
		cobra.CheckErr(err)

		if downloadAll {
			var totalNeeded int64
			for _, site := range sites {
				totalNeeded += site.TotalSize()
			}

			err = client.EnsureFreeSpace(context.Background(), totalNeeded)
			cobra.CheckErr(err)

			for _, site := range sites {
				downloadSite(client, site)
			}
			return
		}

		if len(args) == 0 && !downloadAll {
			siteNames := []string{}
			for _, s := range sites {
				siteNames = append(siteNames, s.Name)
			}

			if len(siteNames) == 0 {
				fmt.Println("No sites found in the last 7 days.")
				return
			}

			selected := []string{}
			prompt := &survey.MultiSelect{
				Message: "Select sites to download:",
				Options: siteNames,
			}
			err = survey.AskOne(prompt, &selected)
			cobra.CheckErr(err)

			selectedSites := []backup.Site{}
			for _, name := range selected {
				for _, site := range sites {
					if site.Name == name {
						selectedSites = append(selectedSites, site)
					}
				}
			}

			var totalNeeded int64
			for _, site := range selectedSites {
				totalNeeded += site.TotalSize()
			}

			err = client.EnsureFreeSpace(context.Background(), totalNeeded)
			cobra.CheckErr(err)

			for _, site := range selectedSites {
				downloadSite(client, site)
			}
			return
		}

		if len(args) > 0 {
			siteName := args[0]
			for _, site := range sites {
				if site.Name == siteName {
					err = client.EnsureFreeSpace(context.Background(), site.TotalSize())
					cobra.CheckErr(err)
					downloadSite(client, site)
					return
				}
			}
			cobra.CheckErr(fmt.Errorf("site %s not found", siteName))
			return
		}

		fmt.Println("Please specify a site or use --all")
	},
}

func downloadSite(client *s3.Client, site backup.Site) {
	fmt.Printf("Site: %s\n", site.Name)
	if site.DB != nil {
		fmt.Printf("  Downloading DB: %s (%s)...\n", site.DB.Key, formatSize(site.DB.Size))
		err := client.DownloadObject(context.Background(), site.DB.Key)
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
		} else {
			fmt.Println("  Done")
		}
	}
	if site.Files != nil {
		fmt.Printf("  Downloading Files: %s (%s)...\n", site.Files.Key, formatSize(site.Files.Size))
		err := client.DownloadObject(context.Background(), site.Files.Key)
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
		} else {
			fmt.Println("  Done")
		}
	}
}

func init() {
	downloadCmd.Flags().BoolVarP(&downloadAll, "all", "a", false, "Download latest backups for all sites")
	rootCmd.AddCommand(downloadCmd)
}
