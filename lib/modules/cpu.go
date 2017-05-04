package modules

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"os"
	"strconv"
	"strings"
)

// Cpu represents the configuration for the cpu block.
type Cpu struct {
	BlockConfigBase `yaml:",inline"`
	Cpu             string  `yaml:"cpu"` // name of the cpu to display (from /proc/stat), empty = all
	CritUsage       float64 `yaml:"crit_usage"`

	lastIdle  uint64
	lastTotal uint64
}

// UpdateBlock updates the status of the CPU block.
// The block displays the CPU usage since the last update.
func (c *Cpu) UpdateBlock(b *i3barjson.Block) {
	b.Color = c.Color

	usage, err := c.getCpuUsage()
	if err != nil {
		b.FullText = err.Error()
		b.Urgent = true
		return
	}

	b.Urgent = (c.CritUsage > 0) && (usage > c.CritUsage)
	b.FullText = fmt.Sprintf("%s%2.0f%%", c.Label, 100*usage)
}

// getCpuUsage returns the usage of the configured CPU. It does by reading /proc/stat and
// calculating (1 - idle_time/total_time).
func (c *Cpu) getCpuUsage() (float64, error) {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return 0, err
	}
	defer f.Close()

	// find the line that corresponds to the configured CPU
	s := bufio.NewScanner(f)
	for s.Scan() {
		flds := strings.Fields(s.Text())
		if c.Cpu == "" || (len(flds) > 0 && flds[0] == c.Cpu) {
			return c.parseCpuLine(flds[1:])
		}
	}
	if s.Err() != nil {
		return 0, s.Err()
	}
	return 0, fmt.Errorf("cpu %s not found", c.Cpu)
}

// parseCpuLine extracts the usage from the text fields of one line of /proc/stat.
func (c *Cpu) parseCpuLine(flds []string) (float64, error) {
	if len(flds) < 7 {
		return 0, errors.New("invalid line in /proc/stat")
	}
	var total, idle uint64
	for i := range flds {
		val, err := strconv.ParseUint(flds[i], 10, 64)
		if err != nil {
			return 0, errors.New("invalid number")
		}
		total += val
		if i == 3 {
			idle += val
		}

	}
	usage := 1 - (float64(idle-c.lastIdle) / float64(total-c.lastTotal))
	c.lastTotal, c.lastIdle = total, idle
	return usage, nil
}
