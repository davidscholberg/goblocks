package main

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"github.com/davidscholberg/goblocks/lib/modules"
	"os"
	"reflect"
)

func main() {
	gb, err := modules.NewGoblocks()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		return
	}

	h := i3barjson.Header{Version: 1}
	err = i3barjson.Init(os.Stdout, nil, h, gb.Cfg.Global.Debug)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		return
	}

	// send the first statusLine
	err = i3barjson.Update(gb.StatusLine)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		return
	}

	for {
		// select on all chans
		i, _, _ := reflect.Select(gb.SelectCases.Cases)
		selectReturn := gb.SelectCases.Actions[i](gb.SelectCases.Blocks[i])
		if selectReturn.Exit {
			break
		}
		if selectReturn.Refresh {
			err = i3barjson.Update(gb.StatusLine)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s", err)
				break
			}
		} else if selectReturn.Reload {
			gb.Reset()
			gb, err = modules.NewGoblocks()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s", err)
				break
			}
			err = i3barjson.Update(gb.StatusLine)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s", err)
				break
			}
		}
	}

	fmt.Println("")
}
