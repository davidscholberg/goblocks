package modules

import (
	"errors"
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

func updateVolumeBlock(b *i3barjson.Block) error {
	amixerCmd := "amixer"
	amixerArgs := []string{"-D", "default", "get", "Master"}
	out, err := exec.Command(amixerCmd, amixerArgs...).Output()
	if err != nil {
		return err
	}
	outStr := string(out)
	iBegin := strings.Index(outStr, "[")
	if iBegin == -1 {
		return errors.New("cannot parse amixer output")
	}
	iEnd := strings.Index(outStr, "]")
	if iEnd == -1 {
		return errors.New("cannot parse amixer output")
	}
	b.FullText = fmt.Sprintf("V: %s", outStr[iBegin+1:iEnd])
	return nil
}
