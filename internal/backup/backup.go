package backup

import (
	"time"
)

type BackupObject struct {
	Key          string
	LastModified time.Time
	Size         int64
	ETag         string
}

type Site struct {
	Name  string
	DB    *BackupObject
	Files *BackupObject
}

func (s *Site) TotalSize() int64 {
	var total int64
	if s.DB != nil {
		total += s.DB.Size
	}
	if s.Files != nil {
		total += s.Files.Size
	}
	return total
}
