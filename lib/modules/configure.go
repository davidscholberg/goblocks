package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"reflect"
	"time"
)

type BlockConfig interface {
	GetBlockIndex() int
	GetUpdateFunc() func(b *i3barjson.Block, c BlockConfig)
	GetUpdateInterval() int
	GetUpdateSignal() int
}

type GoBlock struct {
	Block  i3barjson.Block
	Config BlockConfig
	Ticker *time.Ticker
	Update func(b *i3barjson.Block, c BlockConfig)
}

type Config struct {
	Disk         Disk          `yaml:"disk"`
	Interfaces   []Interface   `yaml:"interfaces"`
	Load         Load          `yaml:"load"`
	Memory       Memory        `yaml:"memory"`
	Raid         Raid          `yaml:"raid"`
	Temperatures []Temperature `yaml:"temperatures"`
	Time         Time          `yaml:"time"`
	Volume       Volume        `yaml:"volume"`
}

func GetGoBlocks(c Config) ([]*GoBlock, error) {
	// TODO: error handling
	// TODO: include i3barjson.Block config in config structs
	var blockConfigs []BlockConfig
	cType := reflect.ValueOf(c)
	for i := 0; i < cType.NumField(); i++ {
		field := cType.Field(i)
		switch field.Kind() {
		case reflect.Struct:
			b := field.Interface().(BlockConfig)
			if b.GetBlockIndex() > 0 {
				blockConfigs = append(blockConfigs, b)
			}
		case reflect.Slice:
			for i := 0; i < field.Len(); i++ {
				b := field.Index(i).Interface().(BlockConfig)
				if b.GetBlockIndex() > 0 {
					blockConfigs = append(blockConfigs, b)
				}
			}
		default:
			return nil, fmt.Errorf("unexpected type: %s\n", field.Type())
		}
	}

	goblocks := make([]*GoBlock, len(blockConfigs))
	for _, blockConfig := range blockConfigs {
		blockIndex := blockConfig.GetBlockIndex()
		updateFunc := blockConfig.GetUpdateFunc()
		ticker := time.NewTicker(time.Second * time.Duration(blockConfig.GetUpdateInterval()))
		goblocks[blockIndex-1] = &GoBlock{
			i3barjson.Block{Separator: true, SeparatorBlockWidth: 20},
			blockConfig,
			ticker,
			updateFunc,
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
