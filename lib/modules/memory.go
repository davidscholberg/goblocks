package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"os"
)

type Memory struct {
	BlockIndex     int     `mapstructure:"block_index"`
	UpdateInterval int     `mapstructure:"update_interval"`
	UpdateSignal   int     `mapstructure:"update_signal"`
	CritMem        float64 `mapstructure:"crit_mem"`
}

func (c Memory) GetBlockIndex() int {
	return c.BlockIndex
}

func (c Memory) GetUpdateFunc() func(b *i3barjson.Block, c BlockConfig) {
	return updateMemBlock
}

func (c Memory) GetUpdateInterval() int {
	return c.UpdateInterval
}

func (c Memory) GetUpdateSignal() int {
	return c.UpdateSignal
}

func updateMemBlock(b *i3barjson.Block, c BlockConfig) {
	cfg := c.(Memory)
	var memAvail, memJunk int64
	fullTextFmt := "M: %s"
	r, err := os.Open("/proc/meminfo")
	if err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	_, err = fmt.Fscanf(
		r,
		"MemTotal: %d kB\nMemFree: %d kB\nMemAvailable: %d ",
		&memJunk, &memJunk, &memAvail)
	if err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	r.Close()
	memAvailG := float64(memAvail) / 1048576.0
	if memAvailG < cfg.CritMem {
		b.Urgent = true
	} else {
		b.Urgent = false
	}
	b.FullText = fmt.Sprintf(fullTextFmt, fmt.Sprintf("%.2fG", memAvailG))
}
