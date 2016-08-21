package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"os"
)

type Interface struct {
	BlockIndex     int    `mapstructure:"block_index"`
	UpdateInterval int    `mapstructure:"update_interval"`
	IfaceName      string `mapstructure:"interface_name"`
}

func (c Interface) GetBlockIndex() int {
	return c.BlockIndex
}

func (c Interface) GetUpdateFunc() func(b *i3barjson.Block, c BlockConfig) {
	return updateIfaceBlock
}

func (c Interface) GetUpdateInterval() int {
	return c.UpdateInterval
}

func updateIfaceBlock(b *i3barjson.Block, c BlockConfig) {
	cfg := c.(Interface)
	var statusStr string
	fullTextFmt := "E: %s"
	// TODO: make interface name configurable
	sysFilePath := fmt.Sprintf("/sys/class/net/%s/operstate", cfg.IfaceName)
	r, err := os.Open(sysFilePath)
	if err != nil {
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	_, err = fmt.Fscanf(r, "%s", &statusStr)
	if err != nil {
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	r.Close()
	b.FullText = fmt.Sprintf(fullTextFmt, statusStr)
}
