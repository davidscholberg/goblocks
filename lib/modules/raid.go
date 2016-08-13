package modules

import (
	"github.com/davidscholberg/go-i3barjson"
	"io/ioutil"
	"strings"
)

func updateRaidBlock(rb *i3barjson.Block) error {
	rb.FullText = "R: ok"
	mdstatPath := "/proc/mdstat"
	stats, err := ioutil.ReadFile(mdstatPath)
	if err != nil {
		return err
	}
	i := strings.Index(string(stats), "_")
	if i != -1 {
		rb.FullText = "R: degraded"
	}
	return nil
}
