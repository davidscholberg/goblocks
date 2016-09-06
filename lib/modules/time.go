package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"time"
)

// Time represents the configuration for the time display block.
type Time struct {
	BlockIndex     int     `yaml:"block_index"`
	UpdateInterval float64 `yaml:"update_interval"`
	Label          string  `yaml:"label"`
	Color          string  `yaml:"color"`
	UpdateSignal   int     `yaml:"update_signal"`
	TimeFormat     string  `yaml:"time_format"`
}

// GetBlockIndex returns the block's position.
func (c Time) GetBlockIndex() int {
	return c.BlockIndex
}

// GetUpdateFunc returns the block's status update function.
func (c Time) GetUpdateFunc() func(b *i3barjson.Block, c BlockConfig) {
	return updateTimeBlock
}

// GetUpdateInterval returns the block's update interval in seconds.
func (c Time) GetUpdateInterval() float64 {
	return c.UpdateInterval
}

// GetUpdateSignal returns the block's update signal that forces an update and
// refresh.
func (c Time) GetUpdateSignal() int {
	return c.UpdateSignal
}

// updateTimeBlock updates the time display block.
func updateTimeBlock(b *i3barjson.Block, c BlockConfig) {
	cfg := c.(Time)
	b.Color = cfg.Color
	b.FullText = fmt.Sprintf(
		"%s%s",
		cfg.Label,
		time.Now().Format(cfg.TimeFormat),
	)
}
