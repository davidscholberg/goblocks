package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"syscall"
)

// Disk represents the configuration for the disk block.
type Disk struct {
	BlockIndex     int     `yaml:"block_index"`
	UpdateInterval float64 `yaml:"update_interval"`
	Label          string  `yaml:"label"`
	UpdateSignal   int     `yaml:"update_signal"`
}

// GetBlockIndex returns the block's position.
func (c Disk) GetBlockIndex() int {
	return c.BlockIndex
}

// GetUpdateFunc returns the block's status update function.
func (c Disk) GetUpdateFunc() func(b *i3barjson.Block, c BlockConfig) {
	return updateDiskBlock
}

// GetUpdateInterval returns the block's update interval in seconds.
func (c Disk) GetUpdateInterval() float64 {
	return c.UpdateInterval
}

// GetUpdateSignal returns the block's update signal that forces an update and
// refresh.
func (c Disk) GetUpdateSignal() int {
	return c.UpdateSignal
}

// updateDiskBlock updates the status of the disk block.
// The block displays "ok" unless one of the given filesystems are over 95%.
func updateDiskBlock(b *i3barjson.Block, c BlockConfig) {
	cfg := c.(Disk)
	fullTextFmt := fmt.Sprintf("%s%%s", cfg.Label)
	fsList := []string{"/", "/home"}
	var err error
	for _, fsPath := range fsList {
		stats := syscall.Statfs_t{}
		err = syscall.Statfs(fsPath, &stats)
		if err != nil {
			b.Urgent = true
			b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
			return
		}
		percentFree := float64(stats.Bavail) * 100 / float64(stats.Blocks)
		if percentFree < 5.0 {
			b.Urgent = true
			b.FullText = fmt.Sprintf(
				fullTextFmt,
				fmt.Sprintf(
					"%s at %.2f%%",
					fsPath,
					100-percentFree,
				),
			)
			return
		}
	}
	b.Urgent = false
	b.FullText = fmt.Sprintf(fullTextFmt, "ok")
}
