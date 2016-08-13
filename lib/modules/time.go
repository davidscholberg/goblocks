package modules

import (
	"github.com/davidscholberg/go-i3barjson"
	"time"
)

func updateTimeBlock(tb *i3barjson.Block) error {
	tb.Full_text = time.Now().Format("2006-01-02 15:04")
	return nil
}
