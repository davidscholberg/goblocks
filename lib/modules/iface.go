package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"os"
)

// Interface represents the configuration for the network interface block.
type Interface struct {
	BlockIndex     int     `yaml:"block_index"`
	UpdateInterval float64 `yaml:"update_interval"`
	Label          string  `yaml:"label"`
	Color          string  `yaml:"color"`
	UpdateSignal   int     `yaml:"update_signal"`
	IfaceName      string  `yaml:"interface_name"`
}

// GetBlockIndex returns the block's position.
func (c Interface) GetBlockIndex() int {
	return c.BlockIndex
}

// GetUpdateFunc returns the block's status update function.
func (c Interface) GetUpdateFunc() func(b *i3barjson.Block, c BlockConfig) {
	return updateIfaceBlock
}

// GetUpdateInterval returns the block's update interval in seconds.
func (c Interface) GetUpdateInterval() float64 {
	return c.UpdateInterval
}

// GetUpdateSignal returns the block's update signal that forces an update and
// refresh.
func (c Interface) GetUpdateSignal() int {
	return c.UpdateSignal
}

// updateIfaceBlock updates the network interface block.
func updateIfaceBlock(b *i3barjson.Block, c BlockConfig) {
	cfg := c.(Interface)
	b.Color = cfg.Color
	fullTextFmt := fmt.Sprintf("%s%%s", cfg.Label)
	var statusStr string
	sysFilePath := fmt.Sprintf("/sys/class/net/%s/operstate", cfg.IfaceName)
	r, err := os.Open(sysFilePath)
	if err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	_, err = fmt.Fscanf(r, "%s", &statusStr)
	if err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	r.Close()
	if statusStr == "up" {
		b.Urgent = false
	} else {
		b.Urgent = true
	}
	b.FullText = fmt.Sprintf(fullTextFmt, statusStr)
}
