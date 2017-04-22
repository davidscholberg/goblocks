package modules

import (
	"fmt"
	"syscall"

	"github.com/davidscholberg/go-i3barjson"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

// Disk represents the configuration for the disk block.
type Disk struct {
	BlockConfigBase `yaml:",inline"`
	Filesystems     map[string]float64 `yaml:"filesystems"`
	DisplayUsage    bool               `yaml:"display_usage"`
}

// UpdateBlock updates the status of the disk block.
// The block displays "ok" unless one of the given filesystems are over 95%.
func (c Disk) UpdateBlock(b *i3barjson.Block) {
	b.Color = c.Color
	stats := syscall.Statfs_t{}
	fullTextFmt := fmt.Sprintf("%s%%s", c.Label)
	if c.DisplayUsage == true {
		diskUsageFmt := ""
		for fsPath, critPercent := range c.Filesystems {
			err := syscall.Statfs(fsPath, &stats)
			if err != nil {
				b.Urgent = true
				b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
				return
			}
			diskFree := float64(stats.Bavail*uint64(stats.Bsize)) / float64(GB)
			diskTotal := float64(stats.Blocks*uint64(stats.Bsize)) / float64(GB)
			b.Urgent = false
			diskUsageFmt += fmt.Sprintf(
				"%s:%.1fG/%.1fG ",
				fsPath,
				float64(diskFree),
				float64(diskTotal),
			)
			percentUsed := 100 - (float64(stats.Bavail) * 100 / float64(stats.Blocks))
			if percentUsed >= critPercent {
				b.Urgent = true
			}

		}
		b.FullText = fmt.Sprintf(fullTextFmt, diskUsageFmt)
	} else {
		for fsPath, critPercent := range c.Filesystems {
			err := syscall.Statfs(fsPath, &stats)
			if err != nil {
				b.Urgent = true
				b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
				return
			}
			percentUsed := 100 - (float64(stats.Bavail) * 100 / float64(stats.Blocks))
			if percentUsed >= critPercent {
				b.FullText = fmt.Sprintf(
					fullTextFmt,
					fmt.Sprintf(
						"%s at %.2f%%",
						fsPath,
						percentUsed,
					),
				)
				b.Urgent = true
				return
			}
		}
		b.Urgent = false
		b.FullText = fmt.Sprintf(fullTextFmt, "ok")
		return
	}
}
