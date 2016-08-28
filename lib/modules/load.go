package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"os"
)

type Load struct {
	BlockIndex     int     `mapstructure:"block_index"`
	UpdateInterval int     `mapstructure:"update_interval"`
	UpdateSignal   int     `mapstructure:"update_signal"`
	CritLoad       float64 `mapstructure:"crit_load"`
}

func (c Load) GetBlockIndex() int {
	return c.BlockIndex
}

func (c Load) GetUpdateFunc() func(b *i3barjson.Block, c BlockConfig) {
	return updateLoadBlock
}

func (c Load) GetUpdateInterval() int {
	return c.UpdateInterval
}

func (c Load) GetUpdateSignal() int {
	return c.UpdateSignal
}

func updateLoadBlock(b *i3barjson.Block, c BlockConfig) {
	cfg := c.(Load)
	var load float64
	fullTextFmt := "L: %.2f"
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
	b.FullText = fmt.Sprintf(fullTextFmt, load)
}
