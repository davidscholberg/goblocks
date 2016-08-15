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

func updateMemBlock(b *i3barjson.Block) {
	var memAvail, memJunk int64
	fullTextFmt := "M: %s"
	r, err := os.Open("/proc/meminfo")
	if err != nil {
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	_, err = fmt.Fscanf(
		r,
		"MemTotal: %d kB\nMemFree: %d kB\nMemAvailable: %d ",
		&memJunk, &memJunk, &memAvail)
	if err != nil {
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	r.Close()
	statusStr := fmt.Sprintf("%.2fG", float64(memAvail)/1048576.0)
	b.FullText = fmt.Sprintf(fullTextFmt, statusStr)
}
