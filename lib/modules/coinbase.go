package modules

import (
	"fmt"

	"github.com/davidscholberg/go-i3barjson"
	"github.com/mnlwldr/coinbase"
)

// Coinbase represents the configuration for the coinbase block.
type Coinbase struct {
	BlockConfigBase `yaml:",inline"`
	Coin            string `yaml:"currency_pair"`
}

// UpdateBlock updates the coinbase block.
// The value dispayed is the current price for the currency_pair
func (c Coinbase) UpdateBlock(b *i3barjson.Block) {
	response, _ := coinbase.Get(c.Coin)
	b.Color = c.Color
	b.FullText = fmt.Sprintf(
		"%s%.8f",
		c.Label,
		response.Amount,
	)
}
