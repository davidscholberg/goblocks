package modules

import (
	"github.com/davidscholberg/go-i3barjson"
	"github.com/davidscholberg/goblocks/lib/types"
	"time"
)

func GetGoBlocks() []*types.GoBlock {
	return []*types.GoBlock{
		getRaidBlock(),
		getDiskBlock(),
		getLoadBlock(),
		getMemBlock(),
		getTempBlock(),
		getIfaceBlock(),
		getVolumeBlock(),
		getTimeBlock(),
	}
}

func newGoBlock(b i3barjson.Block, t *time.Ticker, u func(b *i3barjson.Block) error) *types.GoBlock {
	return &types.GoBlock{b, t, u}
}
