package modules

import (
	"fmt"
	"os"
	"time"

	"github.com/davidscholberg/go-i3barjson"
)

// Uptime represents the configuration for the time display block.
type Uptime struct {
	BlockConfigBase `yaml:",inline"`
}

// UpdateBlock updates the time display block.
func (c Uptime) UpdateBlock(b *i3barjson.Block) {
	b.Color = c.Color
	fullTextFmt := fmt.Sprintf("%s%%s", c.Label)
	var load float64
	r, err := os.Open("/proc/uptime")
	if err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	defer r.Close()
	_, err = fmt.Fscanf(r, "%f ", &load)
	if err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	b.FullText = fmt.Sprintf(
		"%s%s",
		c.Label,
		time.Duration(load)*time.Second,
	)
}
