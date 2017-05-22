package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"os/exec"
	"strings"
)

// Volume represents the configuration for the volume display block.
type Volume struct {
	BlockConfigBase `yaml:",inline"`
	MixerDevice     string `yaml:"mixer_device"`
	Channel         string `yaml:"channel"`
	MuteIndicator   string `yaml:"mute_indicator"`
}

// UpdateBlock updates the volume display block.
// Currently, only the ALSA master channel volume is supported.
func (c Volume) UpdateBlock(b *i3barjson.Block) {
	b.Color = c.Color
	fullTextFmt := fmt.Sprintf("%s%%s", c.Label)
	amixerCmd := "amixer"
	if c.MixerDevice == "" {
		c.MixerDevice = "default"
	}
	if c.Channel == "" {
		c.Channel = "Master"
	}
	if c.MuteIndicator == "" {
		c.MuteIndicator = "muted"
	}
	amixerArgs := []string{"-D", c.MixerDevice, "get", c.Channel}
	out, err := exec.Command(amixerCmd, amixerArgs...).Output()
	if err != nil {
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	outSplit := strings.Split(string(out), "[")
	if len(outSplit) < 3 {
		b.FullText = fmt.Sprintf(fullTextFmt, "cannot parse amixer output")
		return
	}
	statusSplit := outSplit[len(outSplit)-1]
	if statusSplit[:len(statusSplit)-2] == "on" {
		b.FullText = fmt.Sprintf(fullTextFmt, outSplit[1][:len(outSplit[1])-2])
	} else {
		b.FullText = fmt.Sprintf("%s%s", c.Label, c.MuteIndicator)
	}
}
