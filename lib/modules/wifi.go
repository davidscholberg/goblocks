package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"io/ioutil"
	"strings"
)

// Wifi represents the configuration for the wifi percent block.
type Wifi struct {
	BlockConfigBase `yaml:",inline"`
	IfaceName       string  `yaml:"interface_name"`
	CritQuality     float64 `yaml:"crit_quality"`
}

// UpdateBlock updates the wifi block's status.
func (c Wifi) UpdateBlock(b *i3barjson.Block) {
	b.Color = c.Color
	fullTextFmt := fmt.Sprintf("%s%%d%%%%", c.Label)
	var wifiSignal float64
	var ignore string
	wifiPath := "/proc/net/wireless"
	wifiStats, err := ioutil.ReadFile(wifiPath)
	if err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	wifiLines := strings.Split(string(wifiStats), "\n")
	for _, wifiLine := range wifiLines {
		if strings.HasPrefix(wifiLine, c.IfaceName) {
			_, err := fmt.Sscan(wifiLine, &ignore, &ignore, &wifiSignal)
			if err != nil {
				b.Urgent = true
				b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
				return
			}
			break
		}
	}
	wifiPercent := (wifiSignal * 100.0) / 70.0
	if wifiPercent > c.CritQuality {
		b.Urgent = false
	} else {
		b.Urgent = true
	}
	b.FullText = fmt.Sprintf(fullTextFmt, int(wifiPercent))
}
