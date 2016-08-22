package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"os"
)

type Load struct {
	BlockIndex     int `mapstructure:"block_index"`
	UpdateInterval int `mapstructure:"update_interval"`
	UpdateSignal   int `mapstructure:"update_signal"`
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
	var load string
	fullTextFmt := "L: %s"
	r, err := os.Open("/proc/loadavg")
	if err != nil {
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	_, err = fmt.Fscanf(r, "%s ", &load)
	if err != nil {
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	r.Close()
	b.FullText = fmt.Sprintf(fullTextFmt, load)
}
