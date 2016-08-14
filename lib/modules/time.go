package modules

import (
	"github.com/davidscholberg/go-i3barjson"
	"github.com/davidscholberg/goblocks/lib/types"
	"time"
)

func getTimeBlock() *types.GoBlock {
	return newGoBlock(
		i3barjson.Block{Separator: true, SeparatorBlockWidth: 20},
		time.NewTicker(time.Second),
		updateTimeBlock,
	)
}

func updateTimeBlock(b *i3barjson.Block) error {
	b.FullText = time.Now().Format("2006-01-02 15:04")
	return nil
}
