package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"time"
)

// Time represents the configuration for the time display block.
type Time struct {
	BlockConfigBase `yaml:",inline"`
	TimeFormat      string `yaml:"time_format"`
}

// UpdateBlock updates the time display block.
func (c Time) UpdateBlock(b *i3barjson.Block) {
	b.Color = c.Color
	b.FullText = fmt.Sprintf(
		"%s%s",
		c.Label,
		time.Now().Format(c.TimeFormat),
	)
}
