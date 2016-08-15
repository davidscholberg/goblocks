package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"github.com/davidscholberg/goblocks/lib/types"
	"os/exec"
	"strings"
	"time"
)

func getVolumeBlock() *types.GoBlock {
	return newGoBlock(
		i3barjson.Block{Separator: true, SeparatorBlockWidth: 20},
		time.NewTicker(time.Second*60),
		updateVolumeBlock,
	)
}

func updateVolumeBlock(b *i3barjson.Block) {
	fullTextFmt := "V: %s"
	amixerCmd := "amixer"
	amixerArgs := []string{"-D", "default", "get", "Master"}
	out, err := exec.Command(amixerCmd, amixerArgs...).Output()
	if err != nil {
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	outStr := string(out)
	iBegin := strings.Index(outStr, "[")
	if iBegin == -1 {
		b.FullText = fmt.Sprintf(fullTextFmt, "cannot parse amixer output")
		return
	}
	iEnd := strings.Index(outStr, "]")
	if iEnd == -1 {
		b.FullText = fmt.Sprintf(fullTextFmt, "cannot parse amixer output")
		return
	}
	b.FullText = fmt.Sprintf(fullTextFmt, outStr[iBegin+1:iEnd])
}

func SelectActionUpdateVolumeBlock(b *types.GoBlock) (bool, bool) {
	// TODO: fix error handling
	updateVolumeBlock(&b.Block)
	return types.SelectActionRefresh(b)
}
