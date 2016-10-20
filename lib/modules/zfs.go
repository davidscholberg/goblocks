package modules

import (
	"bufio"
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"os/exec"
	"strings"
)

// Zfs represents the configuration data for the ZFS block
type Zfs struct {
	BlockConfigBase `yaml:",inline"`
	PoolName        string `yaml:"zpool_name"`
}

// UpdateBlock updates the ZFS block
func (c Zfs) UpdateBlock(b *i3barjson.Block) {
	b.Color = c.Color
	fullTextFmt := fmt.Sprintf("%s%%s", c.Label)

	zpoolCmd := exec.Command("sudo", "zpool", "status", c.PoolName)
	out, err := zpoolCmd.StdoutPipe()

	if err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}

	if err := zpoolCmd.Start(); err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}

	defer zpoolCmd.Wait()

	buff := bufio.NewScanner(out)
	for buff.Scan() {
		line := strings.TrimSpace(buff.Text())
		if strings.HasPrefix(line, "state") {
			split := strings.Split(line, ":")
			status := strings.TrimSpace(split[1])

			if status == "ONLINE" {
				b.Urgent = false
			} else {
				b.Urgent = true
			}
			b.FullText = fmt.Sprintf(fullTextFmt, status)
			return
		}
	}

	b.Urgent = true
	b.FullText = fmt.Sprintf(fullTextFmt, "NOT FOUND")
	return
}
