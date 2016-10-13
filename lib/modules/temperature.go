package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"io/ioutil"
	"os"
	"strings"
)

// Temperature represents the configuration for the CPU temperature block.
// CpuTempPath is the path to the "hwmon" directory of the CPU temperature info.
// e.g. /sys/devices/platform/coretemp.0/hwmon
type Temperature struct {
	BlockConfigBase `yaml:",inline"`
	CpuTempPath     string  `yaml:"cpu_temp_path"`
	CritTemp        float64 `yaml:"crit_temp"`
}

// UpdateBlock updates the CPU temperature info.
// The value output by the block is the average temperature of all cores.
func (c Temperature) UpdateBlock(b *i3barjson.Block) {
	b.Color = c.Color
	fullTextFmt := fmt.Sprintf("%s%%s", c.Label)
	totalTemp := 0
	procs := 0
	sysFileDirList, err := ioutil.ReadDir(c.CpuTempPath)
	if err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	if len(sysFileDirList) != 1 {
		b.Urgent = true
		msg := fmt.Sprintf(
			"in %s, expected 1 file, got %d",
			c.CpuTempPath,
			len(sysFileDirList),
		)
		b.FullText = fmt.Sprintf(fullTextFmt, msg)
		return
	}
	sysFileDirPath := fmt.Sprintf(
		"%s/%s",
		c.CpuTempPath,
		sysFileDirList[0].Name(),
	)
	sysFileNameFmt := fmt.Sprintf("%s/%%s", sysFileDirPath)
	sysFiles, err := ioutil.ReadDir(sysFileDirPath)
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
		defer r.Close()
		var temp int
		_, err = fmt.Fscanf(r, "%d", &temp)
		if err != nil {
			b.Urgent = true
			b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
			return
		}
		totalTemp += temp
		procs++
	}
	avgTemp := float64(totalTemp) / float64(procs*1000)
	if avgTemp >= c.CritTemp {
		b.Urgent = true
	} else {
		b.Urgent = false
	}
	b.FullText = fmt.Sprintf("%s%.2f°C", c.Label, avgTemp)
}
