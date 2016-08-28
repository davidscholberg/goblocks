package main

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"github.com/davidscholberg/goblocks/lib/modules"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"
)

func main() {
	// TODO: set up default values
	confPath := fmt.Sprintf(
		"%s/.config/goblocks/goblocks.yml",
		os.Getenv("HOME"),
	)
	confStr, err := ioutil.ReadFile(confPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		return
	}

	var cfg modules.Config
	err = yaml.Unmarshal(confStr, &cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		return
	}

	var SIGRTMIN = syscall.Signal(34)

	var statusLine i3barjson.StatusLine
	goblocks, err := modules.GetGoBlocks(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		return
	}

	var selectCases modules.SelectCases
	selectCases.AddBlockSelectCases(goblocks)

	updateTicker := time.NewTicker(time.Second)
	selectCases.AddChanSelectCase(
		updateTicker.C,
		modules.SelectActionRefresh,
	)

	sigEndChan := make(chan os.Signal, 1)
	signal.Notify(sigEndChan, syscall.SIGINT, syscall.SIGTERM)
	selectCases.AddChanSelectCase(
		sigEndChan,
		modules.SelectActionExit,
	)

	for _, goblock := range goblocks {
		statusLine = append(statusLine, &goblock.Block)

		updateSignal := goblock.Config.GetUpdateSignal()
		if updateSignal > 0 {
			sigUpdateChan := make(chan os.Signal, 1)
			signal.Notify(sigUpdateChan, SIGRTMIN+syscall.Signal(updateSignal))
			updateFunc := goblock.Update
			selectCases.Add(
				sigUpdateChan,
				func(b *modules.GoBlock) (bool, bool) {
					updateFunc(&b.Block, b.Config)
					return modules.SelectActionRefresh(b)
				},
				goblock,
			)

		}

		// update block so it's ready for first run
		goblock.Update(&goblock.Block, goblock.Config)
	}

	h := i3barjson.Header{Version: 1}
	err = i3barjson.Init(os.Stdout, nil, h)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		return
	}

	// send the first statusline
	err = i3barjson.Update(statusLine)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		return
	}

	for {
		// select on all chans
		i, _, _ := reflect.Select(selectCases.Cases)
		refresh, exit := selectCases.Actions[i](selectCases.Blocks[i])
		if exit {
			break
		}
		if refresh {
			err = i3barjson.Update(statusLine)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s", err)
				break
			}
		}
	}

	for _, goblock := range goblocks {
		goblock.Ticker.Stop()
	}
	updateTicker.Stop()

	fmt.Println("\ndone")
}
