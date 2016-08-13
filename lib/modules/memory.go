package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"os"
)

func updateMemBlock(mb *i3barjson.Block) error {
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
	mb.FullText = fmt.Sprintf("M: %.2fG", float64(memAvail)/1048576.0)
	return nil
}
