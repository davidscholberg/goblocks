package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"github.com/davidscholberg/goblocks/lib/types"
	"io/ioutil"
	"strings"
	"time"
)

func getRaidBlock() *types.GoBlock {
	return newGoBlock(
		i3barjson.Block{Separator: true, SeparatorBlockWidth: 20},
		time.NewTicker(time.Second),
		updateRaidBlock,
	)
}

func updateRaidBlock(b *i3barjson.Block) {
	fullTextFmt := "R: %s"
	mdstatPath := "/proc/mdstat"
	stats, err := ioutil.ReadFile(mdstatPath)
	if err != nil {
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	i := strings.Index(string(stats), "_")
	if i != -1 {
		b.FullText = fmt.Sprintf(fullTextFmt, "degraded")
		return
	}
	b.FullText = fmt.Sprintf(fullTextFmt, "ok")
}
