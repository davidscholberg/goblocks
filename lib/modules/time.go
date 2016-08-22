package modules

import (
	"github.com/davidscholberg/go-i3barjson"
	"time"
)

type Time struct {
	BlockIndex     int `mapstructure:"block_index"`
	UpdateInterval int `mapstructure:"update_interval"`
	UpdateSignal   int `mapstructure:"update_signal"`
}

func (c Time) GetBlockIndex() int {
	return c.BlockIndex
}

func (c Time) GetUpdateFunc() func(b *i3barjson.Block, c BlockConfig) {
	return updateTimeBlock
}

func (c Time) GetUpdateInterval() int {
	return c.UpdateInterval
}

func (c Time) GetUpdateSignal() int {
	return c.UpdateSignal
}

func updateTimeBlock(b *i3barjson.Block, c BlockConfig) {
	b.FullText = time.Now().Format("2006-01-02 15:04")
}
