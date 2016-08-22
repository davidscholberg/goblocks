package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"syscall"
)

type Disk struct {
	BlockIndex     int `mapstructure:"block_index"`
	UpdateInterval int `mapstructure:"update_interval"`
	UpdateSignal   int `mapstructure:"update_signal"`
}

func (c Disk) GetBlockIndex() int {
	return c.BlockIndex
}

func (c Disk) GetUpdateFunc() func(b *i3barjson.Block, c BlockConfig) {
	return updateDiskBlock
}

func (c Disk) GetUpdateInterval() int {
	return c.UpdateInterval
}

func (c Disk) GetUpdateSignal() int {
	return c.UpdateSignal
}

func updateDiskBlock(b *i3barjson.Block, c BlockConfig) {
	fullTextFmt := "D: %s"
	fsList := []string{"/", "/home"}
	var err error
	for _, fsPath := range fsList {
		stats := syscall.Statfs_t{}
		err = syscall.Statfs(fsPath, &stats)
		if err != nil {
			b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
			return
		}
		percentFree := float64(stats.Bavail) * 100 / float64(stats.Blocks)
		if percentFree < 5.0 {
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
	b.FullText = fmt.Sprintf(fullTextFmt, "ok")
}
