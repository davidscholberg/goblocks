package modules

import (
	"github.com/davidscholberg/go-i3barjson"
	"github.com/davidscholberg/goblocks/lib/types"
	"time"
)

func RegisterGoBlocks(r func(gb []*types.GoBlock)) {
	goblocks := []*types.GoBlock{
		&types.GoBlock{
			&i3barjson.Block{Separator: true, SeparatorBlockWidth: 20},
			time.NewTicker(time.Second),
			updateRaidBlock,
		},
		&types.GoBlock{
			&i3barjson.Block{Separator: true, SeparatorBlockWidth: 20},
			time.NewTicker(time.Second),
			updateDiskBlock,
		},
		&types.GoBlock{
			&i3barjson.Block{Separator: true, SeparatorBlockWidth: 20},
			time.NewTicker(time.Second),
			updateLoadBlock,
		},
		&types.GoBlock{
			&i3barjson.Block{Separator: true, SeparatorBlockWidth: 20},
			time.NewTicker(time.Second),
			updateMemBlock,
		},
		&types.GoBlock{
			&i3barjson.Block{Separator: true, SeparatorBlockWidth: 20},
			time.NewTicker(time.Second),
			updateTempBlock,
		},
		&types.GoBlock{
			&i3barjson.Block{Separator: true, SeparatorBlockWidth: 20},
			time.NewTicker(time.Second),
			updateIfaceBlock,
		},
		&types.GoBlock{
			&i3barjson.Block{Separator: true, SeparatorBlockWidth: 20},
			time.NewTicker(time.Second * 60),
			updateVolumeBlock,
		},
		&types.GoBlock{
			&i3barjson.Block{},
			time.NewTicker(time.Second),
			updateTimeBlock,
		},
	}
	r(goblocks)
}
