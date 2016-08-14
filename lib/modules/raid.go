package modules

import (
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

func updateRaidBlock(b *i3barjson.Block) error {
	b.FullText = "R: ok"
	mdstatPath := "/proc/mdstat"
	stats, err := ioutil.ReadFile(mdstatPath)
	if err != nil {
		return err
	}
	i := strings.Index(string(stats), "_")
	if i != -1 {
		b.FullText = "R: degraded"
	}
	return nil
}
