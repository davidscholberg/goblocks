package modules

import (
	"fmt"
	"os/exec"
	"strings"

	i3barjson "github.com/davidscholberg/go-i3barjson"
)

type Command struct {
	BlockConfigBase `yaml:",inline"`
	Cmd             string `yaml:"command"`
}

func (c Command) UpdateBlock(b *i3barjson.Block) {
	b.Color = c.Color
	fullTextFmt := fmt.Sprintf("%s%%s", c.Label)

	cmdParse := strings.Fields(c.Cmd)
	cmd, args := cmdParse[0], cmdParse[1:]

	if out, err := exec.Command(cmd, args...).Output(); err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}

	b.FullText = fmt.Sprintf(fullTextFmt, strings.Replace(string(out), "\n", "", -1))
}
