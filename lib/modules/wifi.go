package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"io/ioutil"
	"strings"
)

// Wifi represents the configuration for the wifi percent block.
type Wifi struct {
	BlockIndex     int     `yaml:"block_index"`
	UpdateInterval float64 `yaml:"update_interval"`
	Label          string  `yaml:"label"`
	UpdateSignal   int     `yaml:"update_signal"`
	IfaceName      string  `yaml:"interface_name"`
	CritQuality    float64 `yaml:"crit_quality"`
}

// GetBlockIndex returns the block's position.
func (c Wifi) GetBlockIndex() int {
	return c.BlockIndex
}

// GetUpdateFunc returns the block's status update function.
func (c Wifi) GetUpdateFunc() func(b *i3barjson.Block, c BlockConfig) {
	return updateWifiBlock
}

// GetUpdateInterval returns the block's update interval in seconds.
func (c Wifi) GetUpdateInterval() float64 {
	return c.UpdateInterval
}

// GetUpdateSignal returns the block's update signal that forces an update and
// refresh.
func (c Wifi) GetUpdateSignal() int {
	return c.UpdateSignal
}

// updateWifiBlock updates the wifi block's status.
func updateWifiBlock(b *i3barjson.Block, c BlockConfig) {
	cfg := c.(Wifi)
	fullTextFmt := fmt.Sprintf("%s%%d%%%%", cfg.Label)
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
		if strings.HasPrefix(wifiLine, cfg.IfaceName) {
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
	if wifiPercent > cfg.CritQuality {
		b.Urgent = false
	} else {
		b.Urgent = true
	}
	b.FullText = fmt.Sprintf(fullTextFmt, int(wifiPercent))
}
