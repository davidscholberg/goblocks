package types

import (
	"reflect"
)

type SelectCases struct {
	Cases   []reflect.SelectCase
	Actions []func(gb *GoBlock) (refresh bool, exit bool)
	Blocks  []*GoBlock
}

func (s *SelectCases) Add(c interface{}, a func(gb *GoBlock) (bool, bool), b *GoBlock) {
	selectCase := reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(c),
	}
	s.Cases = append(s.Cases, selectCase)
	s.Actions = append(s.Actions, a)
	s.Blocks = append(s.Blocks, b)
}
