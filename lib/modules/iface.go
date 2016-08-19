package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"github.com/davidscholberg/goblocks/lib/types"
	"os"
	"time"
)

var ifaceName string

func getIfaceBlock() *types.GoBlock {
	return newGoBlock(
		i3barjson.Block{Separator: true, SeparatorBlockWidth: 20},
		time.NewTicker(time.Second),
		updateIfaceBlock,
	)
}

func updateIfaceBlock(b *i3barjson.Block) {
	var statusStr string
	fullTextFmt := "E: %s"
	// TODO: make interface name configurable
	sysFilePath := fmt.Sprintf("/sys/class/net/%s/operstate", ifaceName)
	r, err := os.Open(sysFilePath)
	if err != nil {
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	_, err = fmt.Fscanf(r, "%s", &statusStr)
	if err != nil {
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	r.Close()
	b.FullText = fmt.Sprintf(fullTextFmt, statusStr)
}
