package main

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"github.com/davidscholberg/goblocks/lib/modules"
	"log"
	"os"
	"reflect"
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

	for {
		// select on all chans
		i, _, _ := reflect.Select(gb.SelectCases.Cases)
		selectReturn := gb.SelectCases.Actions[i](gb.SelectCases.Blocks[i])
		if selectReturn.Exit {
			fmt.Println("")
			break
		}
		if selectReturn.Refresh {
			err = i3barjson.Update(gb.StatusLine)
			if err != nil {
				errLogger.Fatalln(err)
			}
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
		}
	}
}
