package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"io/ioutil"
	"strings"
)

type Raid struct {
	BlockIndex     int `mapstructure:"block_index"`
	UpdateInterval int `mapstructure:"update_interval"`
}

func (c Raid) GetBlockIndex() int {
	return c.BlockIndex
}

func (c Raid) GetUpdateFunc() func(b *i3barjson.Block, c BlockConfig) {
	return updateRaidBlock
}

func (c Raid) GetUpdateInterval() int {
	return c.UpdateInterval
}

func updateRaidBlock(b *i3barjson.Block, c BlockConfig) {
	fullTextFmt := "R: %s"
	mdstatPath := "/proc/mdstat"
	stats, err := ioutil.ReadFile(mdstatPath)
	if err != nil {
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	i := strings.Index(string(stats), "_")
	if i != -1 {
		b.FullText = fmt.Sprintf(fullTextFmt, "degraded")
		return
	}
	b.FullText = fmt.Sprintf(fullTextFmt, "ok")
}
