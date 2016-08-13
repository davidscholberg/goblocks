package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"io/ioutil"
	"os"
	"strings"
)

func updateTempBlock(tb *i3barjson.Block) error {
	totalTemp := 0
	procs := 0
	sysDirName := "/sys/devices/platform/coretemp.0/hwmon/hwmon1"
	sysFileNameFmt := fmt.Sprintf("%s/%%s", sysDirName)
	sysFiles, err := ioutil.ReadDir(sysDirName)
	if err != nil {
		return err
	}
	for _, sysFile := range sysFiles {
		sysFileName := sysFile.Name()
		if !strings.HasSuffix(sysFileName, "input") {
			continue
		}
		r, err := os.Open(fmt.Sprintf(sysFileNameFmt, sysFileName))
		if err != nil {
			return err
		}
		var temp int
		_, err = fmt.Fscanf(r, "%d", &temp)
		if err != nil {
			return err
		}
		r.Close()
		totalTemp += temp
		procs++
	}
	tb.Full_text = fmt.Sprintf("%.2f°C", float64(totalTemp)/float64(procs*1000))
	return nil
}
