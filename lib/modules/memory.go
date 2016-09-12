package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"os"
)

// Memory represents the configuration for the memory block.
type Memory struct {
	BlockConfigBase `yaml:",inline"`
	CritMem         float64 `yaml:"crit_mem"`
}

// UpdateBlock updates the system memory block status.
// The value dispayed is the amount of available memory.
func (c Memory) UpdateBlock(b *i3barjson.Block) {
	b.Color = c.Color
	fullTextFmt := fmt.Sprintf("%s%%s", c.Label)
	var memAvail, memJunk int64
	r, err := os.Open("/proc/meminfo")
	if err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	defer r.Close()
	_, err = fmt.Fscanf(
		r,
		"MemTotal: %d kB\nMemFree: %d kB\nMemAvailable: %d ",
		&memJunk, &memJunk, &memAvail)
	if err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	memAvailG := float64(memAvail) / 1048576.0
	if memAvailG < c.CritMem {
		b.Urgent = true
	} else {
		b.Urgent = false
	}
	b.FullText = fmt.Sprintf(fullTextFmt, fmt.Sprintf("%.2fG", memAvailG))
}
