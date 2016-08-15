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
	goblocks := modules.GetGoBlocks()
	for _, goblock := range goblocks {
		statusLine = append(statusLine, &goblock.Block)

		// update block so it's ready for first run
		goblock.Update(&goblock.Block)
	}

	var selectCases types.SelectCases
	selectCases.AddBlockSelectCases(goblocks)

	updateTicker := time.NewTicker(time.Second)
	selectCases.AddChanSelectCase(
		updateTicker.C,
		types.SelectActionRefresh,
	)

	sigEndChan := make(chan os.Signal, 1)
	signal.Notify(sigEndChan, syscall.SIGINT, syscall.SIGTERM)
	selectCases.AddChanSelectCase(
		sigEndChan,
		types.SelectActionExit,
	)

	sigVolChan := make(chan os.Signal, 1)
	signal.Notify(sigVolChan, SIGRTMIN+8)
	selectCases.Add(
		sigVolChan,
		modules.SelectActionUpdateVolumeBlock,
		// TODO: don't hardcode this!!!
		goblocks[6],
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
