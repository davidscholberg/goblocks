package modules

import (
	"fmt"
	"github.com/davidscholberg/go-durationfmt"
	"github.com/davidscholberg/go-i3barjson"
	"os"
	"time"
)

// Uptime represents the configuration for the time display block.
type Uptime struct {
	BlockConfigBase `yaml:",inline"`
	DurationFormat  string `yaml:"duration_format"`
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
	dur := time.Duration(load) * time.Second
	durFmt := c.DurationFormat
	if durFmt == "" {
		durFmt = "%hh%mm%ss"
	}
	durStr, err := durationfmt.Format(dur, durFmt)
	if err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	b.FullText = fmt.Sprintf(
		"%s%s",
		c.Label,
		durStr,
	)
}
