package main

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"github.com/davidscholberg/goblocks/lib/modules"
	"github.com/davidscholberg/goblocks/lib/types"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"
)

func main() {
	var SIGRTMIN = syscall.Signal(34)

	var statusLine i3barjson.StatusLine
	var goblocks []*types.GoBlock
	var selectCases types.SelectCases
	modules.RegisterGoBlocks(func(gb []*types.GoBlock) {
		goblocks = gb
		for _, goblock := range goblocks {
			statusLine = append(statusLine, &goblock.Block)
			update := goblock.Update
			selectCases.Add(
				goblock.Ticker.C,
				func(gb *types.GoBlock) (bool, bool) {
					update(&gb.Block)
					return false, false
				},
				goblock,
			)

			// update block so it's ready for first run
			err := goblock.Update(&goblock.Block)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s", err)
			}
		}
	})

	updateTicker := time.NewTicker(time.Second)
	selectCases.Add(
		updateTicker.C,
		func(gb *types.GoBlock) (bool, bool) {
			return true, false
		},
		nil,
	)

	sigEndChan := make(chan os.Signal, 1)
	signal.Notify(sigEndChan, syscall.SIGINT, syscall.SIGTERM)
	selectCases.Add(
		sigEndChan,
		func(gb *types.GoBlock) (bool, bool) {
			return false, true
		},
		nil,
	)

	sigVolChan := make(chan os.Signal, 1)
	signal.Notify(sigVolChan, SIGRTMIN+8)
	// TODO: fix so that we can add the proper block pointer or update function here
	selectCases.Add(
		sigVolChan,
		func(gb *types.GoBlock) (bool, bool) {
			return true, false
		},
		nil,
	)

	h := i3barjson.Header{Version: 1}
	err := i3barjson.Init(os.Stdout, nil, h)
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
