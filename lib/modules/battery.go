package modules

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/davidscholberg/go-i3barjson"
)

// Battery represents the configuration for the battery block.
type Battery struct {
	BlockConfigBase `yaml:",inline"`
	BatteryNumber   int     `yaml:"battery_number"`
	CritBattery     float64 `yaml:"crit_battery"`
	ChargingLabel   string  `yaml:"charging_label"`
}

// UpdateBlock updates the battery status block.
func (c Battery) UpdateBlock(b *i3barjson.Block) {
	var capacity int
	var fullTextFmt string
	b.Color = c.Color

	sysFilePath := fmt.Sprintf("/sys/class/power_supply/BAT%d/capacity", c.BatteryNumber)
	batFilePath := fmt.Sprintf("/sys/class/power_supply/BAT%d/status", c.BatteryNumber)

	r, err := os.Open(batFilePath)
	if err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	defer r.Close()

	batStatus, err := ioutil.ReadAll(r)
	if err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}

	if strings.Contains(string(batStatus), "Charging") {
		fullTextFmt = fmt.Sprintf("%s%%d%%%%", c.ChargingLabel)
	} else {
		fullTextFmt = fmt.Sprintf("%s%%d%%%%", c.Label)
	}

	r, err = os.Open(sysFilePath)
	if err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	defer r.Close()

	_, err = fmt.Fscanf(r, "%d", &capacity)
	if err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	if float64(capacity) >= c.CritBattery {
		b.Urgent = false
	} else {
		b.Urgent = true
	}

	b.FullText = fmt.Sprintf(fullTextFmt, capacity)
}
