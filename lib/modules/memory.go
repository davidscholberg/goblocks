package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"os"
)

// Memory represents the configuration for the memory block.
type Memory struct {
	BlockIndex     int     `yaml:"block_index"`
	UpdateInterval int     `yaml:"update_interval"`
	Label          string  `yaml:"label"`
	UpdateSignal   int     `yaml:"update_signal"`
	CritMem        float64 `yaml:"crit_mem"`
}

// GetBlockIndex returns the block's position.
func (c Memory) GetBlockIndex() int {
	return c.BlockIndex
}

// GetUpdateFunc returns the block's status update function.
func (c Memory) GetUpdateFunc() func(b *i3barjson.Block, c BlockConfig) {
	return updateMemBlock
}

// GetUpdateInterval returns the block's update interval in seconds.
func (c Memory) GetUpdateInterval() int {
	return c.UpdateInterval
}

// GetUpdateSignal returns the block's update signal that forces an update and
// refresh.
func (c Memory) GetUpdateSignal() int {
	return c.UpdateSignal
}

// updateMemBlock updates the system memory block status.
// The value dispayed is the amount of available memory.
func updateMemBlock(b *i3barjson.Block, c BlockConfig) {
	cfg := c.(Memory)
	fullTextFmt := fmt.Sprintf("%s%%s", cfg.Label)
	var memAvail, memJunk int64
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
