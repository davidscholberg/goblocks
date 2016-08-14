package modules

import (
	"github.com/davidscholberg/goblocks/lib/types"
)

func GetBlockSelectCases(b []*types.GoBlock) *types.SelectCases {
	var selectCases types.SelectCases
	for _, goblock := range b {
		addBlockToSelectCase(&selectCases, goblock)
	}
	return &selectCases
}

func addBlockToSelectCase(s *types.SelectCases, b *types.GoBlock) {
	updateFunc := b.Update
	s.Add(
		b.Ticker.C,
		func(gb *types.GoBlock) (bool, bool) {
			updateFunc(&gb.Block)
			return false, false
		},
		b,
	)
}
