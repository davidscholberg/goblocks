package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"io/ioutil"
	"os"
	"strings"
)

// Temperature represents the configuration for the CPU temperature block.
type Temperature struct {
	BlockIndex     int     `yaml:"block_index"`
	UpdateInterval float64 `yaml:"update_interval"`
	Label          string  `yaml:"label"`
	UpdateSignal   int     `yaml:"update_signal"`
	CpuTempPath    string  `yaml:"cpu_temp_path"`
	CritTemp       float64 `yaml:"crit_temp"`
}

// GetBlockIndex returns the block's position.
func (c Temperature) GetBlockIndex() int {
	return c.BlockIndex
}

// GetUpdateFunc returns the block's status update function.
func (c Temperature) GetUpdateFunc() func(b *i3barjson.Block, c BlockConfig) {
	return updateTempBlock
}

// GetUpdateInterval returns the block's update interval in seconds.
func (c Temperature) GetUpdateInterval() float64 {
	return c.UpdateInterval
}

// GetUpdateSignal returns the block's update signal that forces an update and
// refresh.
func (c Temperature) GetUpdateSignal() int {
	return c.UpdateSignal
}

// updateTempBlock updates the CPU temperature info.
// The value output by the block is the average temperature of all cores.
func updateTempBlock(b *i3barjson.Block, c BlockConfig) {
	cfg := c.(Temperature)
	fullTextFmt := fmt.Sprintf("%s%%s", cfg.Label)
	totalTemp := 0
	procs := 0
	sysFileNameFmt := fmt.Sprintf("%s/%%s", cfg.CpuTempPath)
	sysFiles, err := ioutil.ReadDir(cfg.CpuTempPath)
	if err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	for _, sysFile := range sysFiles {
		sysFileName := sysFile.Name()
		if !strings.HasSuffix(sysFileName, "input") {
			continue
		}
		r, err := os.Open(fmt.Sprintf(sysFileNameFmt, sysFileName))
		if err != nil {
			b.Urgent = true
			b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
			return
		}
		var temp int
		_, err = fmt.Fscanf(r, "%d", &temp)
		if err != nil {
			b.Urgent = true
			b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
			return
		}
		r.Close()
		totalTemp += temp
		procs++
	}
	avgTemp := float64(totalTemp) / float64(procs*1000)
	if avgTemp >= cfg.CritTemp {
		b.Urgent = true
	} else {
		b.Urgent = false
	}
	b.FullText = fmt.Sprintf("%s%.2f°C", cfg.Label, avgTemp)
}
