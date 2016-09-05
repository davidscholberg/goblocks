package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"os"
)

// Batter represents the configuration for the battery block.
type Battery struct {
	BlockIndex     int     `yaml:"block_index"`
	UpdateInterval float64 `yaml:"update_interval"`
	Label          string  `yaml:"label"`
	UpdateSignal   int     `yaml:"update_signal"`
	BatteryNumber  int     `yaml:"battery_number"`
	CritBattery    float64 `yaml:"crit_battery"`
}

// GetBlockIndex returns the block's position.
func (c Battery) GetBlockIndex() int {
	return c.BlockIndex
}

// GetUpdateFunc returns the block's status update function.
func (c Battery) GetUpdateFunc() func(b *i3barjson.Block, c BlockConfig) {
	return updateBatteryBlock
}

// GetUpdateInterval returns the block's update interval in seconds.
func (c Battery) GetUpdateInterval() float64 {
	return c.UpdateInterval
}

// GetUpdateSignal returns the block's update signal that forces an update and
// refresh.
func (c Battery) GetUpdateSignal() int {
	return c.UpdateSignal
}

// updateBatteryBlock updates the battery status block.
func updateBatteryBlock(b *i3barjson.Block, c BlockConfig) {
	cfg := c.(Battery)
	fullTextFmt := fmt.Sprintf("%s%%d%%%%", cfg.Label)
	var capacity int
	sysFilePath := fmt.Sprintf("/sys/class/power_supply/BAT%d/capacity", cfg.BatteryNumber)
	r, err := os.Open(sysFilePath)
	if err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	_, err = fmt.Fscanf(r, "%d", &capacity)
	if err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	r.Close()
	if float64(capacity) >= cfg.CritBattery {
		b.Urgent = false
	} else {
		b.Urgent = true
	}
	b.FullText = fmt.Sprintf(fullTextFmt, capacity)
}
