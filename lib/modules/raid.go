package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"io/ioutil"
	"strings"
)

// Raid represents the configuration for the RAID block.
type Raid struct {
	BlockIndex     int     `yaml:"block_index"`
	UpdateInterval float64 `yaml:"update_interval"`
	Label          string  `yaml:"label"`
	Color          string  `yaml:"color"`
	UpdateSignal   int     `yaml:"update_signal"`
}

// GetBlockIndex returns the block's position.
func (c Raid) GetBlockIndex() int {
	return c.BlockIndex
}

// GetUpdateFunc returns the block's status update function.
func (c Raid) GetUpdateFunc() func(b *i3barjson.Block, c BlockConfig) {
	return updateRaidBlock
}

// GetUpdateInterval returns the block's update interval in seconds.
func (c Raid) GetUpdateInterval() float64 {
	return c.UpdateInterval
}

// GetUpdateSignal returns the block's update signal that forces an update and
// refresh.
func (c Raid) GetUpdateSignal() int {
	return c.UpdateSignal
}

// updateRaidBlock updates the RAID block's status.
// This block only supports linux mdraid, and alerts if any RAID volume on the
// system is degraded.
func updateRaidBlock(b *i3barjson.Block, c BlockConfig) {
	cfg := c.(Raid)
	b.Color = cfg.Color
	fullTextFmt := fmt.Sprintf("%s%%s", cfg.Label)
	mdstatPath := "/proc/mdstat"
	stats, err := ioutil.ReadFile(mdstatPath)
	if err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	i := strings.Index(string(stats), "_")
	if i != -1 {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, "degraded")
		return
	}
	b.Urgent = false
	b.FullText = fmt.Sprintf(fullTextFmt, "ok")
}
