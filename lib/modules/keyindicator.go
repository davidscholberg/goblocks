package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"os/exec"
	"strings"
)

// KeyIndicator represents the configuration for the key indicator block.
type KeyIndicator struct {
	BlockConfigBase `yaml:",inline"`
	Keys            []string `yaml:"keys"`
}

// UpdateBlock updates the key indicator block's status.
func (c KeyIndicator) UpdateBlock(b *i3barjson.Block) {
	b.Color = c.Color
	fullTextFmt := fmt.Sprintf("%s%%s", c.Label)
	xsetCmd := "xset"
	xsetArgs := []string{"q"}
	out, err := exec.Command(xsetCmd, xsetArgs...).Output()
	if err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}

	keyStatuses := make(map[string]bool)
	xsetLines := strings.Split(string(out), "\n")
	for _, xsetLine := range xsetLines {
		for _, keyStr := range c.Keys {
			if _, ok := keyStatuses[keyStr]; ok {
				continue
			}
			keyIndex := strings.Index(xsetLine, keyStr)
			if keyIndex == -1 {
				continue
			}
			xsetLineSubstr := xsetLine[keyIndex+len(keyStr):]
			keyStatus := strings.Index(xsetLineSubstr, "o")
			if keyStatus == -1 {
				b.Urgent = true
				b.FullText = fmt.Sprintf(
					fullTextFmt,
					fmt.Sprintf(
						"couldn't find status for key '%s'",
						keyStr,
					),
				)
				return
			}
			switch xsetLineSubstr[keyStatus+1 : keyStatus+2] {
			case "n":
				keyStatuses[keyStr] = true
			case "f":
				keyStatuses[keyStr] = false
			default:
				b.Urgent = true
				b.FullText = fmt.Sprintf(
					fullTextFmt,
					fmt.Sprintf(
						"unknown status for key '%s'",
						keyStr,
					),
				)
				return
			}
		}
		if len(keyStatuses) == len(c.Keys) {
			break
		}
	}
	var keysOn []string
	for keyStr, keyStatus := range keyStatuses {
		if keyStatus {
			keysOn = append(keysOn, keyStr)
		}
	}
	var keysOnStr string
	if len(keysOn) > 0 {
		b.Urgent = true
		keysOnStr = strings.Join(keysOn, ", ")
	} else {
		b.Urgent = false
		keysOnStr = "none"
	}
	b.FullText = fmt.Sprintf(fullTextFmt, keysOnStr)
}
