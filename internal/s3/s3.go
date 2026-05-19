package s3

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/niklas/s3backup/internal/backup"
	cfgpkg "github.com/niklas/s3backup/internal/config"
	"github.com/niklas/s3backup/internal/utils"
)

type Client struct {
	s3Client *s3.Client
	cfg      *cfgpkg.Config
}

func NewClient(cfg *cfgpkg.Config) (*Client, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if cfg.Endpoint != "" {
			return aws.Endpoint{
				URL:           cfg.Endpoint,
				SigningRegion: cfg.Region,
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return &Client{
		s3Client: s3Client,
		cfg:      cfg,
	}, nil
}

func (c *Client) EnsureFreeSpace(ctx context.Context, neededSize int64) error {
	absTargetPath, err := filepath.Abs(c.cfg.TargetPath)
	if err != nil {
		return err
	}

	disk, err := utils.GetDiskSpace(absTargetPath)
	if err != nil {
		// If directory doesn't exist, try its parent
		disk, err = utils.GetDiskSpace(filepath.Dir(absTargetPath))
		if err != nil {
			return err
		}
	}

	minFree, err := utils.ParseSize(c.cfg.KeepFree, disk.All)
	if err != nil {
		return err
	}

	if disk.Free >= uint64(neededSize)+minFree {
		return nil
	}

	// Not enough space, need to delete old files in TargetPath
	var files []struct {
		path string
		info os.FileInfo
	}

	err = filepath.Walk(absTargetPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}
		if !info.IsDir() {
			files = append(files, struct {
				path string
				info os.FileInfo
			}{path, info})
		}
		return nil
	})
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].info.ModTime().Before(files[j].info.ModTime())
	})

	currentFree := disk.Free
	for _, f := range files {
		if currentFree >= uint64(neededSize)+minFree {
			break
		}
		fmt.Printf("  Freeing up space: deleting %s\n", f.path)
		err := os.Remove(f.path)
		if err == nil {
			currentFree += uint64(f.info.Size())
		}
	}

	if currentFree < uint64(neededSize)+minFree {
		return fmt.Errorf("could not free up enough space (needed: %d, min_free: %d, available after cleanup: %d)", neededSize, minFree, currentFree)
	}

	return nil
}

func (c *Client) DownloadObject(ctx context.Context, key string) error {
	input := &s3.GetObjectInput{
		Bucket: aws.String(c.cfg.Bucket),
		Key:    aws.String(key),
	}

	result, err := c.s3Client.GetObject(ctx, input)
	if err != nil {
		return err
	}
	defer result.Body.Close()

	target := filepath.Join(c.cfg.TargetPath, key)
	if _, err := os.Stat(target); err == nil {
		// File already exists
		return nil
	}

	err = os.MkdirAll(filepath.Dir(target), 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(target)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, result.Body)
	return err
}

func (c *Client) GetLatestBackups(ctx context.Context) ([]backup.Site, error) {
	var sites []backup.Site
	siteMap := make(map[string]*backup.Site)

	paginator := s3.NewListObjectsV2Paginator(c.s3Client, &s3.ListObjectsV2Input{
		Bucket: aws.String(c.cfg.Bucket),
	})

	maxAge := time.Duration(c.cfg.MaxAge) * 24 * time.Hour
	now := time.Now()

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, obj := range page.Contents {
			if aws.ToInt64(obj.Size) == 0 {
				continue
			}

			if now.Sub(aws.ToTime(obj.LastModified)) > maxAge {
				continue
			}

			key := aws.ToString(obj.Key)
			parts := strings.Split(key, "/")
			if len(parts) < 2 {
				continue
			}

			siteName := parts[0]
			fileName := parts[1]

			s, ok := siteMap[siteName]
			if !ok {
				s = &backup.Site{Name: siteName}
				siteMap[siteName] = s
			}

			isDB := strings.HasSuffix(fileName, ".sql.gz")
			backupObj := &backup.BackupObject{
				Key:          key,
				LastModified: aws.ToTime(obj.LastModified),
				Size:         aws.ToInt64(obj.Size),
				ETag:         aws.ToString(obj.ETag),
			}

			if isDB {
				if s.DB == nil || backupObj.LastModified.After(s.DB.LastModified) {
					s.DB = backupObj
				}
			} else {
				if s.Files == nil || backupObj.LastModified.After(s.Files.LastModified) {
					s.Files = backupObj
				}
			}
		}
	}

	for _, s := range siteMap {
		sites = append(sites, *s)
	}

	return sites, nil
}
