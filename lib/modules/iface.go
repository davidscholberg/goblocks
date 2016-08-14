package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"github.com/davidscholberg/goblocks/lib/types"
	"os"
	"time"
)

func getIfaceBlock() *types.GoBlock {
	return newGoBlock(
		i3barjson.Block{Separator: true, SeparatorBlockWidth: 20},
		time.NewTicker(time.Second),
		updateIfaceBlock,
	)
}

func updateIfaceBlock(b *i3barjson.Block) error {
	var ifaceState string
	// TODO: make interface name configurable
	r, err := os.Open("/sys/class/net/enp3s0/operstate")
	if err != nil {
		return err
	}
	_, err = fmt.Fscanf(r, "%s", &ifaceState)
	if err != nil {
		return err
	}
	r.Close()
	b.FullText = fmt.Sprintf("E: %s", ifaceState)
	return nil
}
