package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"net"
	"strings"
)

// Interface represents the configuration for the network interface block.
type Interface struct {
	BlockConfigBase `yaml:",inline"`
	IfaceName       string `yaml:"interface_name"`
	IPv4            string
	IPv4CIDR        string
	IPv6            string
	IPv6CIDR        string
}

// UpdateBlock updates the network interface block.
func (c Interface) UpdateBlock(b *i3barjson.Block) {
	b.Color = c.Color
	fullTextFmt := fmt.Sprintf("%s%%s", c.Label)

	iface, err := net.InterfaceByName(c.IfaceName)
	if err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}

	if iface.Flags == net.FlagUp {
		b.Urgent = false
	} else {
		b.Urgent = true
	}

	addrs, err := iface.Addrs()
	if err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}

	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			b.Urgent = true
			b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
			return
		}

		// Checking for address family
		if ip.To4() != nil {
			c.Label = strings.Replace(c.Label, "\u003cipv4\u003e", ip.String(), -1)
			c.Label = strings.Replace(c.Label, "\u003ccidr4\u003e", addr.String(), -1)
		} else {
			c.Label = strings.Replace(c.Label, "\u003cipv6\u003e", ip.String(), -1)
			c.Label = strings.Replace(c.Label, "\u003ccidr6\u003e", addr.String(), -1)
		}
	}

	b.FullText = fmt.Sprintf(fullTextFmt, c.Label)
}
