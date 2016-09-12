package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"os"
)

// Interface represents the configuration for the network interface block.
type Interface struct {
	BlockConfigBase `yaml:",inline"`
	IfaceName       string `yaml:"interface_name"`
}

// UpdateBlock updates the network interface block.
func (c Interface) UpdateBlock(b *i3barjson.Block) {
	b.Color = c.Color
	fullTextFmt := fmt.Sprintf("%s%%s", c.Label)
	var statusStr string
	sysFilePath := fmt.Sprintf("/sys/class/net/%s/operstate", c.IfaceName)
	r, err := os.Open(sysFilePath)
	if err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	defer r.Close()
	_, err = fmt.Fscanf(r, "%s", &statusStr)
	if err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	if statusStr == "up" {
		b.Urgent = false
	} else {
		b.Urgent = true
	}
	b.FullText = fmt.Sprintf(fullTextFmt, statusStr)
}
