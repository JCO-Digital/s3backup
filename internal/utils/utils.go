package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

type DiskSpace struct {
	All  uint64
	Free uint64
}

func GetDiskSpace(path string) (DiskSpace, error) {
	var stat syscall.Statfs_t
	err := syscall.Statfs(path, &stat)
	if err != nil {
		return DiskSpace{}, err
	}

	all := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bavail * uint64(stat.Bsize)

	return DiskSpace{All: all, Free: free}, nil
}

func ParseSize(sizeStr string, totalDisk uint64) (uint64, error) {
	sizeStr = strings.TrimSpace(strings.ToUpper(sizeStr))

	if strings.HasSuffix(sizeStr, "%") {
		percStr := strings.TrimSuffix(sizeStr, "%")
		perc, err := strconv.ParseFloat(percStr, 64)
		if err != nil {
			return 0, err
		}
		return uint64(float64(totalDisk) * (perc / 100.0)), nil
	}

	re := regexp.MustCompile(`^(\d+(?:\.\d+)?)\s*([KMGT]?)B?$`)
	matches := re.FindStringSubmatch(sizeStr)
	if matches == nil {
		return 0, fmt.Errorf("invalid size format: %s", sizeStr)
	}

	val, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, err
	}

	var multiplier uint64 = 1
	switch matches[2] {
	case "T":
		multiplier = 1024 * 1024 * 1024 * 1024
	case "G":
		multiplier = 1024 * 1024 * 1024
	case "M":
		multiplier = 1024 * 1024
	case "K":
		multiplier = 1024
	}

	return uint64(val * float64(multiplier)), nil
}
