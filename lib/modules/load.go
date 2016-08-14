package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"github.com/davidscholberg/goblocks/lib/types"
	"os"
	"time"
)

func getLoadBlock() *types.GoBlock {
	return newGoBlock(
		i3barjson.Block{Separator: true, SeparatorBlockWidth: 20},
		time.NewTicker(time.Second),
		updateLoadBlock,
	)
}

func updateLoadBlock(b *i3barjson.Block) error {
	var load string
	r, err := os.Open("/proc/loadavg")
	if err != nil {
		return err
	}
	_, err = fmt.Fscanf(r, "%s ", &load)
	if err != nil {
		return err
	}
	r.Close()
	b.FullText = fmt.Sprintf("L: %s", load)
	return nil
}
