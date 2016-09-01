package main

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"github.com/davidscholberg/goblocks/lib/modules"
	"os"
	"reflect"
)

func main() {
	var cfg modules.Config
	err := modules.GetConfig(&cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		return
	}

	goblocks, err := modules.GetGoBlocks(cfg.Blocks)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		return
	}

	var selectCases modules.SelectCases
	selectCases.AddBlockSelectCases(goblocks)
	selectCases.AddSignalSelectCases(goblocks)
	selectCases.AddUpdateTickerSelectCase(cfg.Global.RefreshInterval)

	h := i3barjson.Header{Version: 1}
	err = i3barjson.Init(os.Stdout, nil, h, cfg.Global.Debug)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		return
	}

	var statusLine i3barjson.StatusLine
	for _, goblock := range goblocks {
		statusLine = append(statusLine, &goblock.Block)
		// update block so it's ready for first run
		goblock.Update(&goblock.Block, goblock.Config)
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
	selectCases.UpdateTicker.Stop()

	fmt.Println("\ndone")
}
