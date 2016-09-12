package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"os"
)

// Load represents the configuration for the system load.
type Load struct {
	BlockConfigBase `yaml:",inline"`
	CritLoad        float64 `yaml:"crit_load"`
}

// UpdateBlock updates the load block status.
func (c Load) UpdateBlock(b *i3barjson.Block) {
	b.Color = c.Color
	fullTextFmt := fmt.Sprintf("%s%%s", c.Label)
	var load float64
	r, err := os.Open("/proc/loadavg")
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
	if load >= c.CritLoad {
		b.Urgent = true
	} else {
		b.Urgent = false
	}
	b.FullText = fmt.Sprintf("%s%.2f", c.Label, load)
}
