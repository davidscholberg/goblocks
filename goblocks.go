package main

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/davidscholberg/go-i3barjson"
	"github.com/davidscholberg/goblocks/lib/modules"
)

func main() {
	errLogger := log.New(os.Stderr, "error: ", 0)

	gb, err := modules.NewGoblocks()
	if err != nil {
		errLogger.Fatalln(err)
	}

	h := i3barjson.Header{Version: 1}
	err = i3barjson.Init(os.Stdout, nil, h, gb.Cfg.Global.Debug)
	if err != nil {
		errLogger.Fatalln(err)
	}

	// send the first statusLine
	err = i3barjson.Update(gb.StatusLine)
	if err != nil {
		errLogger.Fatalln(err)
	}

	shouldRefresh := false

	for {
		// select on all chans
		i, _, _ := reflect.Select(gb.SelectCases.Cases)
		selectReturn := gb.SelectCases.Actions[i](gb.SelectCases.Blocks[i])
		if selectReturn.Exit {
			fmt.Println("")
			break
		}
		if selectReturn.SignalRefresh {
			shouldRefresh = true
		} else if selectReturn.Refresh {
			if shouldRefresh {
				err = i3barjson.Update(gb.StatusLine)
				if err != nil {
					errLogger.Fatalln(err)
				}
				shouldRefresh = false
			}
		} else if selectReturn.ForceRefresh {
			err = i3barjson.Update(gb.StatusLine)
			if err != nil {
				errLogger.Fatalln(err)
			}
			shouldRefresh = false
		} else if selectReturn.Reload {
			gb.Reset()
			gb, err = modules.NewGoblocks()
			if err != nil {
				errLogger.Fatalln(err)
			}
			err = i3barjson.Update(gb.StatusLine)
			if err != nil {
				errLogger.Fatalln(err)
			}
			shouldRefresh = false
		}
	}
}
