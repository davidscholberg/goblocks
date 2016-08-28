package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"io/ioutil"
	"os"
	"strings"
)

type Temperature struct {
	BlockIndex     int     `yaml:"block_index"`
	UpdateInterval int     `yaml:"update_interval"`
	UpdateSignal   int     `yaml:"update_signal"`
	CpuTempPath    string  `yaml:"cpu_temp_path"`
	CritTemp       float64 `yaml:"crit_temp"`
}

func (c Temperature) GetBlockIndex() int {
	return c.BlockIndex
}

func (c Temperature) GetUpdateFunc() func(b *i3barjson.Block, c BlockConfig) {
	return updateTempBlock
}

func (c Temperature) GetUpdateInterval() int {
	return c.UpdateInterval
}

func (c Temperature) GetUpdateSignal() int {
	return c.UpdateSignal
}

func updateTempBlock(b *i3barjson.Block, c BlockConfig) {
	cfg := c.(Temperature)
	totalTemp := 0
	procs := 0
	sysFileNameFmt := fmt.Sprintf("%s/%%s", cfg.CpuTempPath)
	sysFiles, err := ioutil.ReadDir(cfg.CpuTempPath)
	if err != nil {
		b.Urgent = true
		b.FullText = err.Error()
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
			b.FullText = err.Error()
			return
		}
		var temp int
		_, err = fmt.Fscanf(r, "%d", &temp)
		if err != nil {
			b.Urgent = true
			b.FullText = err.Error()
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
	b.FullText = fmt.Sprintf("%.2f°C", avgTemp)
}
