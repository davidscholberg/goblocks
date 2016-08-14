package modules

import (
	"github.com/davidscholberg/go-i3barjson"
	"github.com/davidscholberg/goblocks/lib/types"
	"time"
)

func RegisterGoBlocks(r func(gb []*types.GoBlock)) {
	goblocks := []*types.GoBlock{
		getRaidBlock(),
		getDiskBlock(),
		getLoadBlock(),
		getMemBlock(),
		getTempBlock(),
		getIfaceBlock(),
		getVolumeBlock(),
		getTimeBlock(),
	}
	r(goblocks)
}

func newGoBlock(b i3barjson.Block, t *time.Ticker, u func(b *i3barjson.Block) error) *types.GoBlock {
	return &types.GoBlock{b, t, u}
}
