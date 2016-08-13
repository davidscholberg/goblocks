package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"os"
)

func updateLoadBlock(lb *i3barjson.Block) error {
	var load string
	r, err := os.Open("/proc/loadavg")
	if err != nil {
		return err
	}
	_, err = fmt.Fscanf(r, "%s ", &load)
	if err != nil {
		return err
	}
	r.Close()
	lb.Full_text = fmt.Sprintf("L: %s", load)
	return nil
}
