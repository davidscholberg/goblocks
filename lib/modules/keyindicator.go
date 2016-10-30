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
	Key             string `yaml:"key"`
	KeyText         string `yaml:"key-text"`
	OnColor         string `yaml:"on-color"`
	OffColor        string `yaml:"off-color"`
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

	keyFound := false
	keyStatus := false
	xsetLines := strings.Split(string(out), "\n")
	for _, xsetLine := range xsetLines {
		keyIndex := strings.Index(xsetLine, c.Key)
		if keyIndex == -1 {
			continue
		}
		xsetLineSubstr := xsetLine[keyIndex+len(c.Key):]
		keyStatusIndex := strings.Index(xsetLineSubstr, "o")
		if keyStatusIndex == -1 {
			b.Urgent = true
			b.FullText = fmt.Sprintf(
				fullTextFmt,
				fmt.Sprintf(
					"couldn't find status for key '%s'",
					c.Key,
				),
			)
			return
		}
		switch xsetLineSubstr[keyStatusIndex+1 : keyStatusIndex+2] {
		case "n":
			keyFound = true
			keyStatus = true
		case "f":
			keyFound = true
			keyStatus = false
		default:
			b.Urgent = true
			b.FullText = fmt.Sprintf(
				fullTextFmt,
				fmt.Sprintf(
					"unknown status for key '%s'",
					c.Key,
				),
			)
			return
		}
		if keyFound {
			break
		}
	}

	if !keyFound {
		b.Urgent = true
		b.FullText = fmt.Sprintf(
			fullTextFmt,
			fmt.Sprintf(
				"couldn't find key '%s'",
				c.Key,
			),
		)
		return
	}

	b.Urgent = false

	if keyStatus {
		b.Color = c.OnColor
	} else {
		b.Color = c.OffColor
	}

	b.FullText = fmt.Sprintf(fullTextFmt, c.KeyText)
}
