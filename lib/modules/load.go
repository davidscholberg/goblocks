package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"os"
)

// Load represents the configuration for the system load.
type Load struct {
	BlockIndex     int     `yaml:"block_index"`
	UpdateInterval float64 `yaml:"update_interval"`
	Label          string  `yaml:"label"`
	UpdateSignal   int     `yaml:"update_signal"`
	CritLoad       float64 `yaml:"crit_load"`
}

// GetBlockIndex returns the block's position.
func (c Load) GetBlockIndex() int {
	return c.BlockIndex
}

// GetUpdateFunc returns the block's status update function.
func (c Load) GetUpdateFunc() func(b *i3barjson.Block, c BlockConfig) {
	return updateLoadBlock
}

// GetUpdateInterval returns the block's update interval in seconds.
func (c Load) GetUpdateInterval() float64 {
	return c.UpdateInterval
}

// GetUpdateSignal returns the block's update signal that forces an update and
// refresh.
func (c Load) GetUpdateSignal() int {
	return c.UpdateSignal
}

// updateLoadBlock updates the load block status.
func updateLoadBlock(b *i3barjson.Block, c BlockConfig) {
	cfg := c.(Load)
	fullTextFmt := fmt.Sprintf("%s%%s", cfg.Label)
	var load float64
	r, err := os.Open("/proc/loadavg")
	if err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	_, err = fmt.Fscanf(r, "%f ", &load)
	if err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	r.Close()
	if load >= cfg.CritLoad {
		b.Urgent = true
	} else {
		b.Urgent = false
	}
	b.FullText = fmt.Sprintf("%s%.2f", cfg.Label, load)
}
