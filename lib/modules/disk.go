package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"syscall"
)

func updateDiskBlock(db *i3barjson.Block) error {
	db.Full_text = "D: ok"
	fsList := []string{"/", "/home"}
	var err error
	for _, fsPath := range fsList {
		stats := syscall.Statfs_t{}
		err = syscall.Statfs(fsPath, &stats)
		if err != nil {
			return err
		}
		percentFree := float64(stats.Bavail) * 100 / float64(stats.Blocks)
		if percentFree < 5.0 {
			db.Full_text = fmt.Sprintf(
				"D: %s at %.2f%%",
				fsPath,
				100-percentFree,
			)
			return nil
		}
	}
	return nil
}
