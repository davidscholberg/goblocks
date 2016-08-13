package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"os"
)

func updateIfaceBlock(ib *i3barjson.Block) error {
	var ifaceState string
	// TODO: make interface name configurable
	r, err := os.Open("/sys/class/net/enp3s0/operstate")
	if err != nil {
		return err
	}
	_, err = fmt.Fscanf(r, "%s", &ifaceState)
	if err != nil {
		return err
	}
	r.Close()
	ib.Full_text = fmt.Sprintf("E: %s", ifaceState)
	return nil
}
