package modules

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	i3barjson "github.com/davidscholberg/go-i3barjson"
)

type Command struct {
	BlockConfigBase `yaml:",inline"`
	Cmd             string `yaml:"command"`
	Append          string `yaml:"append"`
	CritValue       string `yaml:"crit_value"`
	CritOperator    string `yaml:"crit_operator"`
}

func (c Command) UpdateBlock(b *i3barjson.Block) {
	var cmdOutput []byte
	b.Color = c.Color
	fullTextFmt := fmt.Sprintf("%s%%s", c.Label)

	cmdParse := strings.Fields(c.Cmd)
	cmd, args := cmdParse[0], cmdParse[1:]

	cmdOutput, err := exec.Command(cmd, args...).Output()
	if err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}

	// init vars as false
	convertToIntegers := false
	b.Urgent = false

	// if we detect a < or > we have to try to convert to ints
	if strings.Contains(c.CritOperator, "<") || strings.Contains(c.CritOperator, ">") {
		convertToIntegers = true
	}

	// trim the script output
	trimmedOutput := strings.TrimSpace(string(cmdOutput))

	// begin int stuff
	if convertToIntegers == true {
		// try to convert script output to int
		intOutput, err := strconv.Atoi(trimmedOutput)
		if err != nil {
			b.Urgent = true
			msg := fmt.Sprintf("script output '%s' is not an int", trimmedOutput)
			b.FullText = fmt.Sprintf(fullTextFmt, msg)
			return
		}

		// try to convert crit_value to int
		intCritValue, err := strconv.Atoi(c.CritValue)
		if err != nil {
			b.Urgent = true
			msg := fmt.Sprintf("crit_value is not an int")
			b.FullText = fmt.Sprintf(fullTextFmt, msg)
			return
		}

		// no errors, safe to do a number comparison
		switch c.CritOperator {
		case ">":
			if intOutput > intCritValue {
				b.Urgent = true
			}
		case "<":
			if intOutput < intCritValue {
				b.Urgent = true
			}
		}
	}

	// safe to check equals either way
	if c.CritOperator == "=" && c.CritValue == trimmedOutput {
		b.Urgent = true
	}

	b.FullText = fmt.Sprintf(fullTextFmt, trimmedOutput+c.Append)
}
