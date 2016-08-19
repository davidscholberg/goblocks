package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"github.com/davidscholberg/goblocks/lib/types"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var sysDirName string

func getTempBlock() *types.GoBlock {
	return newGoBlock(
		i3barjson.Block{Separator: true, SeparatorBlockWidth: 20},
		time.NewTicker(time.Second),
		updateTempBlock,
	)
}

func updateTempBlock(b *i3barjson.Block) {
	totalTemp := 0
	procs := 0
	sysFileNameFmt := fmt.Sprintf("%s/%%s", sysDirName)
	sysFiles, err := ioutil.ReadDir(sysDirName)
	if err != nil {
		b.FullText = err.Error()
		return
	}
	for _, sysFile := range sysFiles {
		sysFileName := sysFile.Name()
		if !strings.HasSuffix(sysFileName, "input") {
			continue
		}
		r, err := os.Open(fmt.Sprintf(sysFileNameFmt, sysFileName))
		if err != nil {
			b.FullText = err.Error()
			return
		}
		var temp int
		_, err = fmt.Fscanf(r, "%d", &temp)
		if err != nil {
			b.FullText = err.Error()
			return
		}
		r.Close()
		totalTemp += temp
		procs++
	}
	b.FullText = fmt.Sprintf("%.2f°C", float64(totalTemp)/float64(procs*1000))
}
