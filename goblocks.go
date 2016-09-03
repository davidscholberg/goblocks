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
	var selectCases modules.SelectCases
	var statusLine i3barjson.StatusLine
	err := modules.Init(&cfg, &selectCases, &statusLine)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		return
	}

	h := i3barjson.Header{Version: 1}
	err = i3barjson.Init(os.Stdout, nil, h, cfg.Global.Debug)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		return
	}

	// send the first statusLine
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

	selectCases.StopTickers()

	fmt.Println("")
}
