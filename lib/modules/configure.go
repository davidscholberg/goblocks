package modules

import (
	"github.com/davidscholberg/go-i3barjson"
	"reflect"
	"time"
)

type BlockConfig interface {
	GetBlockIndex() int
	GetUpdateFunc() func(b *i3barjson.Block, c BlockConfig)
	GetUpdateInterval() int
}

type GoBlock struct {
	Block  i3barjson.Block
	Config BlockConfig
	Ticker *time.Ticker
	Update func(b *i3barjson.Block, c BlockConfig)
}

type Config struct {
	Disk        Disk        `mapstructure:"disk"`
	Interface   Interface   `mapstructure:"interface"`
	Load        Load        `mapstructure:"load"`
	Memory      Memory      `mapstructure:"memory"`
	Raid        Raid        `mapstructure:"raid"`
	Temperature Temperature `mapstructure:"temperature"`
	Time        Time        `mapstructure:"time"`
	Volume      Volume      `mapstructure:"volume"`
}

func GetGoBlocks(c Config) ([]*GoBlock, error) {
	// TODO: error handling
	// TODO: include i3barjson.Block config in config structs
	cType := reflect.ValueOf(c)
	goblocksSize := 0
	for i := 0; i < cType.NumField(); i++ {
		blockConfig := cType.Field(i).Interface().(BlockConfig)
		blockIndex := blockConfig.GetBlockIndex()
		if blockIndex > 0 {
			goblocksSize++
		}
	}

	goblocks := make([]*GoBlock, goblocksSize)
	for i := 0; i < cType.NumField(); i++ {
		blockConfig := cType.Field(i).Interface().(BlockConfig)
		blockIndex := blockConfig.GetBlockIndex()
		if blockIndex > 0 {
			updateFunc := blockConfig.GetUpdateFunc()
			ticker := time.NewTicker(time.Second * time.Duration(blockConfig.GetUpdateInterval()))
			goblocks[blockIndex-1] = &GoBlock{
				i3barjson.Block{Separator: true, SeparatorBlockWidth: 20},
				blockConfig,
				ticker,
				updateFunc,
			}
		}
	}

	return goblocks, nil
}

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
			updateFunc(&b.Block, b.Config)
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
