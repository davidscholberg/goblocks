package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"syscall"
)

// Disk represents the configuration for the disk block.
type Disk struct {
	BlockConfigBase `yaml:",inline"`
	Filesystems     map[string]float64 `yaml:"filesystems"`
}

// UpdateBlock updates the status of the disk block.
// The block displays "ok" unless one of the given filesystems are over 95%.
func (c Disk) UpdateBlock(b *i3barjson.Block) {
	b.Color = c.Color
	fullTextFmt := fmt.Sprintf("%s%%s", c.Label)
	for fsPath, critPercent := range c.Filesystems {
		stats := syscall.Statfs_t{}
		err := syscall.Statfs(fsPath, &stats)
		if err != nil {
			b.Urgent = true
			b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
			return
		}
		percentUsed := 100 - (float64(stats.Bavail) * 100 / float64(stats.Blocks))
		if percentUsed >= critPercent {
			b.Urgent = true
			b.FullText = fmt.Sprintf(
				fullTextFmt,
				fmt.Sprintf(
					"%s at %.2f%%",
					fsPath,
					percentUsed,
				),
			)
			return
		}
	}
	b.Urgent = false
	b.FullText = fmt.Sprintf(fullTextFmt, "ok")
}
