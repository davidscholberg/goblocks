package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"io/ioutil"
	"strings"
)

// Raid represents the configuration for the RAID block.
type Raid struct {
	BlockConfigBase `yaml:",inline"`
}

// UpdateBlock updates the RAID block's status.
// This block only supports linux mdraid, and alerts if any RAID volume on the
// system is degraded.
func (c Raid) UpdateBlock(b *i3barjson.Block) {
	b.Color = c.Color
	fullTextFmt := fmt.Sprintf("%s%%s", c.Label)
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
