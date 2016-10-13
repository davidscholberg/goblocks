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

// Block contains all functions and objects necessary to configure and update
// a single status block.
type Block struct {
	I3barBlock i3barjson.Block
	Config     BlockConfig
}

// PreConfig is the struct used to initially unmarshal the configuration. Once
// the configuration has been fully processed, it is stored in the Config
// struct.
type PreConfig struct {
	Global GlobalConfig             `yaml:"global"`
	Blocks []map[string]interface{} `yaml:"blocks"`
}

// Config is the root configuration struct.
type Config struct {
	Global GlobalConfig
	Blocks []BlockConfig
}

// GlobalConfig represents global config options.
type GlobalConfig struct {
	Debug           bool    `yaml:"debug"`
	RefreshInterval float64 `yaml:"refresh_interval"`
}

// BlockConfig is an interface for Block configuration structs.
type BlockConfig interface {
	GetUpdateInterval() float64
	GetUpdateSignal() int
	GetBlockType() string
	UpdateBlock(b *i3barjson.Block)
}

// BlockConfigBase is a base struct for Block configuration structs. It
// implements all of the methods of the BlockConfig interface except the
// UpdateBlock method. That method should be implemented by each Block
// configuration struct, which should also embed the BlockConfigBase struct as
// an anonymous field. That way, each Block configuration struct will implement
// the full BlockConfig interface.
type BlockConfigBase struct {
	Type           string  `yaml:"type"`
	UpdateInterval float64 `yaml:"update_interval"`
	Label          string  `yaml:"label"`
	Color          string  `yaml:"color"`
	UpdateSignal   int     `yaml:"update_signal"`
}

// GetUpdateInterval returns the block's update interval in seconds.
func (c BlockConfigBase) GetUpdateInterval() float64 {
	return c.UpdateInterval
}

// GetUpdateSignal returns the block's update signal that forces an update and
// refresh.
func (c BlockConfigBase) GetUpdateSignal() int {
	return c.UpdateSignal
}

// GetBlockType returns the block's type as a string.
func (c BlockConfigBase) GetBlockType() string {
	return c.Type
}

// getBlockConfigInstance returns a BlockConfig object whose underlying type is
// determined from the passed-in config map.
func getBlockConfigInstance(m map[string]interface{}) (*BlockConfig, error) {
	yamlStr, err := yaml.Marshal(m)
	if err != nil {
		return nil, err
	}
	t := m["type"].(string)
	switch t {
	case "battery":
		c := Battery{}
		err := yaml.Unmarshal(yamlStr, &c)
		b := BlockConfig(c)
		return &b, err
	case "disk":
		c := Disk{}
		err := yaml.Unmarshal(yamlStr, &c)
		b := BlockConfig(c)
		return &b, err
	case "interface":
		c := Interface{}
		err := yaml.Unmarshal(yamlStr, &c)
		b := BlockConfig(c)
		return &b, err
	case "load":
		c := Load{}
		err := yaml.Unmarshal(yamlStr, &c)
		b := BlockConfig(c)
		return &b, err
	case "memory":
		c := Memory{}
		err := yaml.Unmarshal(yamlStr, &c)
		b := BlockConfig(c)
		return &b, err
	case "raid":
		c := Raid{}
		err := yaml.Unmarshal(yamlStr, &c)
		b := BlockConfig(c)
		return &b, err
	case "temperature":
		c := Temperature{}
		err := yaml.Unmarshal(yamlStr, &c)
		b := BlockConfig(c)
		return &b, err
	case "time":
		c := Time{}
		err := yaml.Unmarshal(yamlStr, &c)
		b := BlockConfig(c)
		return &b, err
	case "volume":
		c := Volume{}
		err := yaml.Unmarshal(yamlStr, &c)
		b := BlockConfig(c)
		return &b, err
	case "wifi":
		c := Wifi{}
		err := yaml.Unmarshal(yamlStr, &c)
		b := BlockConfig(c)
		return &b, err
	}
	return nil, fmt.Errorf("error: type %s not valid", t)
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

	preCfg := PreConfig{}
	err = yaml.Unmarshal(confStr, &preCfg)
	if err != nil {
		return err
	}

	cfg.Global = preCfg.Global

	for _, m := range preCfg.Blocks {
		block, err := getBlockConfigInstance(m)
		if err != nil {
			return err
		}
		cfg.Blocks = append(cfg.Blocks, *block)
	}

	return nil
}

// GetBlocks initializes and returns a Block slice based on the
// given configuration.
func GetBlocks(blockConfigSlice []BlockConfig) ([]*Block, error) {
	blocks := make([]*Block, len(blockConfigSlice))
	for i, blockConfig := range blockConfigSlice {
		blocks[i] = &Block{
			i3barjson.Block{Separator: true, SeparatorBlockWidth: 20},
			blockConfig,
		}
	}

	return blocks, nil
}

// SelectCases represents the set of channels that Goblocks selects on in the
// main program loop, as well as the functions and data to run and operate on,
// respectively.
type SelectCases struct {
	Cases   []reflect.SelectCase
	Actions []SelectAction
	Blocks  []*Block
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
			s.add(
				sigUpdateChan,
				func(b *Block) SelectReturn {
					b.Config.UpdateBlock(&b.I3barBlock)
					return SelectActionRefresh(b)
				},
				block,
			)

		}
	}
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
func addBlockToSelectCase(s *SelectCases, b *Block, c <-chan time.Time) {
	s.add(
		c,
		func(b *Block) SelectReturn {
			b.Config.UpdateBlock(&b.I3barBlock)
			return SelectActionNoop(b)
		},
		b,
	)
}

// SelectAction is a function type that performs an action when a channel is
// selected on in the main program loop. The return value indicates some action
// for the caller to take.
type SelectAction func(*Block) SelectReturn

// SelectReturn is returned by a SelectAction type function and tells the caller
// a certain action to take.
type SelectReturn struct {
	Refresh bool
	Reload  bool
	Exit    bool
}

// SelectActionExit is a helper function of type SelectAction that tells
// Goblocks to exit.
func SelectActionExit(b *Block) SelectReturn {
	return SelectReturn{Exit: true}
}

// SelectActionRefresh is a helper function of type SelectAction that tells
// Goblocks to refresh the output.
func SelectActionRefresh(b *Block) SelectReturn {
	return SelectReturn{Refresh: true}
}

// SelectActionReload is a helper function of type SelectAction that tells
// Goblocks to reload the configuration.
func SelectActionReload(b *Block) SelectReturn {
	return SelectReturn{Reload: true}
}

// SelectActionNoop is a helper function of type SelectAction that tells
// Goblocks to not perform any SelectReturn actions.
func SelectActionNoop(b *Block) SelectReturn {
	return SelectReturn{}
}

// Goblocks contains all configuration and runtime data needed for the
// application.
type Goblocks struct {
	Cfg         Config
	SelectCases SelectCases
	Tickers     []*time.Ticker
	StatusLine  i3barjson.StatusLine
}

// NewGoblocks returns a Goblocks instance with all configuration and runtime
// data loaded in.
func NewGoblocks() (*Goblocks, error) {
	gb := Goblocks{}
	err := GetConfig(&gb.Cfg)
	if err != nil {
		return nil, err
	}

	blocks, err := GetBlocks(gb.Cfg.Blocks)
	if err != nil {
		return nil, err
	}

	gb.SelectCases.AddSignalSelectCases(blocks)
	gb.AddBlockSelectCases(blocks)
	gb.AddUpdateTickerSelectCase()

	for _, block := range blocks {
		gb.StatusLine = append(gb.StatusLine, &block.I3barBlock)
		// update block so it's ready for first run
		block.Config.UpdateBlock(&block.I3barBlock)
	}

	return &gb, nil
}

// AddBlockSelectCases is a helper function to add all configured Block
// objects to Goblocks' SelectCases.
func (gb *Goblocks) AddBlockSelectCases(b []*Block) {
	for _, block := range b {
		ticker := time.NewTicker(
			time.Duration(
				block.Config.GetUpdateInterval() * float64(time.Second),
			),
		)
		gb.Tickers = append(gb.Tickers, ticker)
		addBlockToSelectCase(&gb.SelectCases, block, ticker.C)
	}
}

// AddUpdateTickerSelectCase adds the Goblocks update ticker that controls
// refreshing the Goblocks output.
func (gb *Goblocks) AddUpdateTickerSelectCase() {
	updateTicker := time.NewTicker(
		time.Duration(gb.Cfg.Global.RefreshInterval * float64(time.Second)),
	)
	gb.SelectCases.addChanSelectCase(
		updateTicker.C,
		SelectActionRefresh,
	)
	gb.Tickers = append(gb.Tickers, updateTicker)
}

// Reset stops all tickers and resets all signal handlers.
func (gb *Goblocks) Reset() {
	for _, ticker := range gb.Tickers {
		ticker.Stop()
	}
	signal.Reset()
}
