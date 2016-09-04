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

// Block contains all functions and objects necessary to configure and update
// a block.
type Block struct {
	I3barBlock i3barjson.Block
	Config     BlockConfig
	Ticker     *time.Ticker
	Update     func(b *i3barjson.Block, c BlockConfig)
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

// GetBlocks initializes and returns a Block slice based on the
// given configuration.
func GetBlocks(c BlockConfigs) ([]*Block, error) {
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

	blocks := make([]*Block, len(blockConfigSlice))
	for _, blockConfig := range blockConfigSlice {
		blockIndex := blockConfig.GetBlockIndex()
		updateFunc := blockConfig.GetUpdateFunc()
		ticker := time.NewTicker(
			time.Duration(
				blockConfig.GetUpdateInterval() * float64(time.Second),
			),
		)
		blocks[blockIndex-1] = &Block{
			i3barjson.Block{Separator: true, SeparatorBlockWidth: 20},
			blockConfig,
			ticker,
			updateFunc,
		}
	}

	return blocks, nil
}

// SelectCases represents the set of channels that Goblocks selects on in the
// main program loop, as well as the functions and data to run and operate on,
// respectively.
type SelectCases struct {
	Cases        []reflect.SelectCase
	Actions      []SelectAction
	Blocks       []*Block
	UpdateTicker *time.Ticker
}

// AddBlockSelectCases is a helper function to add all configured Block
// objects to SelectCases.
func (s *SelectCases) AddBlockSelectCases(b []*Block) {
	for _, block := range b {
		addBlockToSelectCase(s, block)
	}
}

const sigrtmin = syscall.Signal(34)

// AddSignalSelectCases loads the select cases related to OS signals.
func (s *SelectCases) AddSignalSelectCases(blocks []*Block) {
	sigReloadChan := make(chan os.Signal, 1)
	signal.Notify(sigReloadChan, syscall.SIGHUP)
	s.addChanSelectCase(
		sigReloadChan,
		SelectActionReload,
	)

	sigEndChan := make(chan os.Signal, 1)
	signal.Notify(sigEndChan, syscall.SIGINT, syscall.SIGTERM)
	s.addChanSelectCase(
		sigEndChan,
		SelectActionExit,
	)

	for _, block := range blocks {
		updateSignal := block.Config.GetUpdateSignal()
		if updateSignal > 0 {
			sigUpdateChan := make(chan os.Signal, 1)
			signal.Notify(sigUpdateChan, sigrtmin+syscall.Signal(updateSignal))
			updateFunc := block.Update
			s.add(
				sigUpdateChan,
				func(b *Block) (bool, bool, bool) {
					updateFunc(&b.I3barBlock, b.Config)
					return SelectActionRefresh(b)
				},
				block,
			)

		}
	}
}

// AddUpdateTickerSelectCase is a helper function to add the update ticker chan
// to SelectCases. The ticker will tick at the given refreshInterval.
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

// add adds a channel, action, and Block to the SelectCases object.
func (s *SelectCases) add(c interface{}, a SelectAction, b *Block) {
	selectCase := reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(c),
	}
	s.Cases = append(s.Cases, selectCase)
	s.Actions = append(s.Actions, a)
	s.Blocks = append(s.Blocks, b)
}

// addChanSelectCase is a helper function that adds a non-Block channel and
// action to SelectCases. This can be used for signal handling and other non-
// block specific operations.
func (s *SelectCases) addChanSelectCase(c interface{}, a SelectAction) {
	s.add(
		c,
		a,
		nil,
	)
}

// addBlockToSelectCase is a helper function to add a Block to SelectCases.
// The channel used is a time.Ticker channel set to tick according to the
// block's configuration. The SelectAction function updates the block's status
// but does not tell Goblocks to refresh.
func addBlockToSelectCase(s *SelectCases, b *Block) {
	updateFunc := b.Update
	s.add(
		b.Ticker.C,
		func(b *Block) (bool, bool, bool) {
			updateFunc(&b.I3barBlock, b.Config)
			return false, false, false
		},
		b,
	)
}

// Reset stops all tickers and resets all signal handlers.
func (s *SelectCases) Reset() {
	for _, block := range s.Blocks {
		if block != nil {
			block.Ticker.Stop()
		}
	}
	s.UpdateTicker.Stop()
	signal.Reset()
}

// Init initializes the configuration, SelectCases, and StatusLine
func Init(cfg *Config, selectCases *SelectCases, statusLine *i3barjson.StatusLine) error {
	err := GetConfig(cfg)
	if err != nil {
		return err
	}

	blocks, err := GetBlocks(cfg.Blocks)
	if err != nil {
		return err
	}

	selectCases.AddBlockSelectCases(blocks)
	selectCases.AddSignalSelectCases(blocks)
	selectCases.AddUpdateTickerSelectCase(cfg.Global.RefreshInterval)

	for _, block := range blocks {
		*statusLine = append(*statusLine, &block.I3barBlock)
		// update block so it's ready for first run
		block.Update(&block.I3barBlock, block.Config)
	}

	return nil
}

// SelectAction is a function type that specifies an action to perform when a
// channel is selected on in the main program loop. The first returned bool
// indicates whether or not Goblocks should refresh the output. The second
// returned bool indicates whether or not to reload the Goblocks configuration.
// The third returned bool indicates whether or not Goblocks should exit the
// loop.
type SelectAction func(*Block) (bool, bool, bool)

// SelectActionExit is a helper function of type SelectAction that tells
// Goblocks to exit.
func SelectActionExit(b *Block) (bool, bool, bool) {
	return false, false, true
}

// SelectActionRefresh is a helper function of type SelectAction that tells
// Goblocks to refresh the output.
func SelectActionRefresh(b *Block) (bool, bool, bool) {
	return true, false, false
}

// SelectActionReload is a helper function of type SelectAction that tells
// Goblocks to reload the configuration.
func SelectActionReload(b *Block) (bool, bool, bool) {
	return false, true, false
}
