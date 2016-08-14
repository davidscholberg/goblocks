package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"github.com/davidscholberg/goblocks/lib/types"
	"os"
	"time"
)

func getMemBlock() *types.GoBlock {
	return newGoBlock(
		i3barjson.Block{Separator: true, SeparatorBlockWidth: 20},
		time.NewTicker(time.Second),
		updateMemBlock,
	)
}

func updateMemBlock(b *i3barjson.Block) error {
	var memAvail, memJunk int64
	r, err := os.Open("/proc/meminfo")
	if err != nil {
		return err
	}
	_, err = fmt.Fscanf(
		r,
		"MemTotal: %d kB\nMemFree: %d kB\nMemAvailable: %d ",
		&memJunk, &memJunk, &memAvail)
	if err != nil {
		return err
	}
	r.Close()
	b.FullText = fmt.Sprintf("M: %.2fG", float64(memAvail)/1048576.0)
	return nil
}
