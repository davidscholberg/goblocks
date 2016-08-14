package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"github.com/davidscholberg/goblocks/lib/types"
	"syscall"
	"time"
)

func getDiskBlock() *types.GoBlock {
	return newGoBlock(
		i3barjson.Block{Separator: true, SeparatorBlockWidth: 20},
		time.NewTicker(time.Second),
		updateDiskBlock,
	)
}

func updateDiskBlock(b *i3barjson.Block) error {
	b.FullText = "D: ok"
	fsList := []string{"/", "/home"}
	var err error
	for _, fsPath := range fsList {
		stats := syscall.Statfs_t{}
		err = syscall.Statfs(fsPath, &stats)
		if err != nil {
			return err
		}
		percentFree := float64(stats.Bavail) * 100 / float64(stats.Blocks)
		if percentFree < 5.0 {
			b.FullText = fmt.Sprintf(
				"D: %s at %.2f%%",
				fsPath,
				100-percentFree,
			)
			return nil
		}
	}
	return nil
}
