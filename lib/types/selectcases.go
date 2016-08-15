package types

import (
	"reflect"
)

type SelectAction func(*GoBlock) (bool, bool)

type SelectCases struct {
	Cases   []reflect.SelectCase
	Actions []SelectAction
	Blocks  []*GoBlock
}

func (s *SelectCases) Add(c interface{}, a SelectAction, b *GoBlock) {
	selectCase := reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(c),
	}
	s.Cases = append(s.Cases, selectCase)
	s.Actions = append(s.Actions, a)
	s.Blocks = append(s.Blocks, b)
}

func (s *SelectCases) AddBlockSelectCases(b []*GoBlock) {
	for _, goblock := range b {
		addBlockToSelectCase(s, goblock)
	}
}

func (s *SelectCases) AddChanSelectCase(c interface{}, a SelectAction) {
	s.Add(
		c,
		a,
		nil,
	)
}

func addBlockToSelectCase(s *SelectCases, b *GoBlock) {
	updateFunc := b.Update
	s.Add(
		b.Ticker.C,
		func(b *GoBlock) (bool, bool) {
			updateFunc(&b.Block)
			return false, false
		},
		b,
	)
}

func SelectActionExit(b *GoBlock) (bool, bool) {
	return false, true
}

func SelectActionRefresh(b *GoBlock) (bool, bool) {
	return true, false
}
