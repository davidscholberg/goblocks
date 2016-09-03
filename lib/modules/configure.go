package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"
)

// BlockConfig is an interface for Block configuration structs.
type BlockConfig interface {
	GetBlockIndex() int
	GetUpdateFunc() func(b *i3barjson.Block, c BlockConfig)
	GetUpdateInterval() float64
	GetUpdateSignal() int
}

// GlobalConfig represents global config options.
type GlobalConfig struct {
	Debug           bool    `yaml:"debug"`
	RefreshInterval float64 `yaml:"refresh_interval"`
}

// GoBlock contains all functions and objects necessary to configure and update
// a block.
type GoBlock struct {
	Block  i3barjson.Block
	Config BlockConfig
	Ticker *time.Ticker
	Update func(b *i3barjson.Block, c BlockConfig)
}

// Config is the root configuration struct.
type Config struct {
	Global GlobalConfig `yaml:"global"`
	Blocks BlockConfigs `yaml:"blocks"`
}

// BlockConfigs holds the configuration of all blocks.
type BlockConfigs struct {
	Disk         Disk          `yaml:"disk"`
	Interfaces   []Interface   `yaml:"interfaces"`
	Load         Load          `yaml:"load"`
	Memory       Memory        `yaml:"memory"`
	Raid         Raid          `yaml:"raid"`
	Temperatures []Temperature `yaml:"temperatures"`
	Time         Time          `yaml:"time"`
	Volume       Volume        `yaml:"volume"`
}

const confPathFmt = "%s/.config/goblocks/goblocks.yml"
const sigrtmin = syscall.Signal(34)

// GetConfig loads the Goblocks configuration object.
func GetConfig(cfg *Config) error {
	// TODO: set up default values
	confPath := fmt.Sprintf(confPathFmt, os.Getenv("HOME"))
	confStr, err := ioutil.ReadFile(confPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(confStr, cfg)
	if err != nil {
		return err
	}

	return nil
}

// GetGoBlocks initializes and returns a GoBlock slice based on the
// given configuration.
func GetGoBlocks(c BlockConfigs) ([]*GoBlock, error) {
	// TODO: error handling
	// TODO: include i3barjson.Block config in config structs
	var blockConfigSlice []BlockConfig
	cType := reflect.ValueOf(c)
	for i := 0; i < cType.NumField(); i++ {
		field := cType.Field(i)
		switch field.Kind() {
		case reflect.Struct:
			b := field.Interface().(BlockConfig)
			if b.GetBlockIndex() > 0 {
				blockConfigSlice = append(blockConfigSlice, b)
			}
		case reflect.Slice:
			for i := 0; i < field.Len(); i++ {
				b := field.Index(i).Interface().(BlockConfig)
				if b.GetBlockIndex() > 0 {
					blockConfigSlice = append(blockConfigSlice, b)
				}
			}
		default:
			return nil, fmt.Errorf("unexpected type: %s\n", field.Type())
		}
	}

	goblocks := make([]*GoBlock, len(blockConfigSlice))
	for _, blockConfig := range blockConfigSlice {
		blockIndex := blockConfig.GetBlockIndex()
		updateFunc := blockConfig.GetUpdateFunc()
		ticker := time.NewTicker(
			time.Duration(
				blockConfig.GetUpdateInterval() * float64(time.Second),
			),
		)
		goblocks[blockIndex-1] = &GoBlock{
			i3barjson.Block{Separator: true, SeparatorBlockWidth: 20},
			blockConfig,
			ticker,
			updateFunc,
		}
	}

	return goblocks, nil
}

// SelectAction is a function type that specifies an action to perform when a
// channel is selected on in the main program loop. The first returned bool
// indicates whether or not Goblocks should refresh the output. The second
// returned bool indicates whether or not Goblocks should exit the loop.
type SelectAction func(*GoBlock) (bool, bool)

// SelectCases represents the set of channels that Goblocks selects on in the
// main program loop, as well as the functions and data to run and operate on,
// respectively.
type SelectCases struct {
	Cases        []reflect.SelectCase
	Actions      []SelectAction
	Blocks       []*GoBlock
	UpdateTicker *time.Ticker
}

// AddBlockSelectCases is a helper function to add all configured GoBlock
// objects to SelectCases.
func (s *SelectCases) AddBlockSelectCases(b []*GoBlock) {
	for _, goblock := range b {
		addBlockToSelectCase(s, goblock)
	}
}

// AddSignalSelectCases loads the select cases related to OS signals.
func (s *SelectCases) AddSignalSelectCases(goblocks []*GoBlock) {
	sigEndChan := make(chan os.Signal, 1)
	signal.Notify(sigEndChan, syscall.SIGINT, syscall.SIGTERM)
	s.addChanSelectCase(
		sigEndChan,
		SelectActionExit,
	)

	for _, goblock := range goblocks {
		updateSignal := goblock.Config.GetUpdateSignal()
		if updateSignal > 0 {
			sigUpdateChan := make(chan os.Signal, 1)
			signal.Notify(sigUpdateChan, sigrtmin+syscall.Signal(updateSignal))
			updateFunc := goblock.Update
			s.add(
				sigUpdateChan,
				func(b *GoBlock) (bool, bool) {
					updateFunc(&b.Block, b.Config)
					return SelectActionRefresh(b)
				},
				goblock,
			)

		}
	}
}

func (s *SelectCases) AddUpdateTickerSelectCase(refreshInterval float64) {
	updateTicker := time.NewTicker(
		time.Duration(refreshInterval * float64(time.Second)),
	)
	s.addChanSelectCase(
		updateTicker.C,
		SelectActionRefresh,
	)
	s.UpdateTicker = updateTicker
}

// add adds a channel, action, and GoBlock to the SelectCases object.
func (s *SelectCases) add(c interface{}, a SelectAction, b *GoBlock) {
	selectCase := reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(c),
	}
	s.Cases = append(s.Cases, selectCase)
	s.Actions = append(s.Actions, a)
	s.Blocks = append(s.Blocks, b)
}

// addChanSelectCase is a helper function that adds a non-GoBlock channel and
// action to SelectCases. This can be used for signal handling and other non-
// block specific operations.
func (s *SelectCases) addChanSelectCase(c interface{}, a SelectAction) {
	s.add(
		c,
		a,
		nil,
	)
}

// addBlockToSelectCase is a helper function to add a GoBlock to SelectCases.
// The channel used is a time.Ticker channel set to tick according to the
// block's configuration. The SelectAction function updates the block's status
// but does not tell Goblocks to refresh.
func addBlockToSelectCase(s *SelectCases, b *GoBlock) {
	updateFunc := b.Update
	s.add(
		b.Ticker.C,
		func(b *GoBlock) (bool, bool) {
			updateFunc(&b.Block, b.Config)
			return false, false
		},
		b,
	)
}

// StopTickers stops all tickers in the SelectCases object.
func (s *SelectCases) StopTickers() {
	for _, goblock := range s.Blocks {
		if goblock != nil {
			goblock.Ticker.Stop()
		}
	}
	s.UpdateTicker.Stop()
}

// Init initializes the configuration, SelectCases, and StatusLine
func Init(cfg *Config, selectCases *SelectCases, statusLine *i3barjson.StatusLine) error {
	err := GetConfig(cfg)
	if err != nil {
		return err
	}

	goblocks, err := GetGoBlocks(cfg.Blocks)
	if err != nil {
		return err
	}

	selectCases.AddBlockSelectCases(goblocks)
	selectCases.AddSignalSelectCases(goblocks)
	selectCases.AddUpdateTickerSelectCase(cfg.Global.RefreshInterval)

	for _, goblock := range goblocks {
		*statusLine = append(*statusLine, &goblock.Block)
		// update block so it's ready for first run
		goblock.Update(&goblock.Block, goblock.Config)
	}

	return nil
}

// SelectActionExit is a helper function of type SelectAction that tells
// Goblocks to exit.
func SelectActionExit(b *GoBlock) (bool, bool) {
	return false, true
}

// SelectActionRefresh is a helper function of type SelectAction that tells
// Goblocks to refresh the output.
func SelectActionRefresh(b *GoBlock) (bool, bool) {
	return true, false
}
