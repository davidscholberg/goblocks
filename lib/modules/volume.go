package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"os/exec"
	"strings"
)

// Volume represents the configuration for the volume display block.
type Volume struct {
	BlockIndex     int     `yaml:"block_index"`
	UpdateInterval float64 `yaml:"update_interval"`
	Label          string  `yaml:"label"`
	UpdateSignal   int     `yaml:"update_signal"`
}

// GetBlockIndex returns the block's position.
func (c Volume) GetBlockIndex() int {
	return c.BlockIndex
}

// GetUpdateFunc returns the block's status update function.
func (c Volume) GetUpdateFunc() func(b *i3barjson.Block, c BlockConfig) {
	return updateVolumeBlock
}

// GetUpdateInterval returns the block's update interval in seconds.
func (c Volume) GetUpdateInterval() float64 {
	return c.UpdateInterval
}

// GetUpdateSignal returns the block's update signal that forces an update and
// refresh.
func (c Volume) GetUpdateSignal() int {
	return c.UpdateSignal
}

// updateVolumeBlock updates the volume display block.
// Currently, only the ALSA master channel volume is supported.
func updateVolumeBlock(b *i3barjson.Block, c BlockConfig) {
	cfg := c.(Volume)
	fullTextFmt := fmt.Sprintf("%s%%s", cfg.Label)
	amixerCmd := "amixer"
	amixerArgs := []string{"-D", "default", "get", "Master"}
	out, err := exec.Command(amixerCmd, amixerArgs...).Output()
	if err != nil {
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	outStr := string(out)
	iBegin := strings.Index(outStr, "[")
	if iBegin == -1 {
		b.FullText = fmt.Sprintf(fullTextFmt, "cannot parse amixer output")
		return
	}
	iEnd := strings.Index(outStr, "]")
	if iEnd == -1 {
		b.FullText = fmt.Sprintf(fullTextFmt, "cannot parse amixer output")
		return
	}
	b.FullText = fmt.Sprintf(fullTextFmt, outStr[iBegin+1:iEnd])
}
