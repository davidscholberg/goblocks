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
	var selectCases []reflect.SelectCase
	modules.RegisterGoBlocks(func(gb []*types.GoBlock) {
		goblocks = gb
		for _, goblock := range goblocks {
			statusLine = append(statusLine, goblock.Block)
			selectCases = append(selectCases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(goblock.Ticker.C),
			})

			// update block so it's ready for first run
			err := goblock.Update(goblock.Block)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s", err)
			}
		}
	})

	updateTicker := time.NewTicker(time.Second)

	selectCases = append(selectCases, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(updateTicker.C),
	})
	updateTickerIndex := len(selectCases) - 1

	sigEndChan := make(chan os.Signal, 1)
	signal.Notify(sigEndChan, syscall.SIGINT, syscall.SIGTERM)

	selectCases = append(selectCases, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(sigEndChan),
	})
	sigEndChanIndex := len(selectCases) - 1

	sigVolChan := make(chan os.Signal, 1)
	signal.Notify(sigVolChan, SIGRTMIN+8)

	selectCases = append(selectCases, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(sigVolChan),
	})
	sigVolChanIndex := len(selectCases) - 1

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
		i, _, _ := reflect.Select(selectCases)
		if i == sigEndChanIndex {
			break
		}
		if i == updateTickerIndex {
			err = i3barjson.Update(statusLine)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s", err)
				break
			}
		} else if i == sigVolChanIndex {
			// TODO: terrible hack, need to reference block by string or var
			err = goblocks[6].Update(goblocks[6].Block)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s", err)
			}
			err = i3barjson.Update(statusLine)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s", err)
				break
			}
		} else {
			err = goblocks[i].Update(goblocks[i].Block)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s", err)
			}
		}
	}

	for _, goblock := range goblocks {
		goblock.Ticker.Stop()
	}
	updateTicker.Stop()

	fmt.Println("\ndone")
}
