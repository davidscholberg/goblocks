package types

import (
	"github.com/davidscholberg/go-i3barjson"
	"time"
)

type GoBlock struct {
	Block  *i3barjson.Block
	Ticker *time.Ticker
	Update func(b *i3barjson.Block) error
}
