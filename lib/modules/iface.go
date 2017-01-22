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
	IPv6Local       string
}

// UpdateBlock updates the network interface block.
func (c Interface) UpdateBlock(b *i3barjson.Block) {
	var (
		status string
	)

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
		status = "up"
	} else {
		b.Urgent = true
		status = "down"
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
			fullTextFmt = strings.Replace(fullTextFmt, "\u003cipv4\u003e", ip.String(), -1)
			fullTextFmt = strings.Replace(fullTextFmt, "\u003ccidr4\u003e", addr.String(), -1)
		} else {

			if ip.String()[0:4] == "fe80" {
				//setting ipv6 link local
				fullTextFmt = strings.Replace(fullTextFmt, "\u003clocal6\u003e", ip.String(), -1)
			} else {
				fullTextFmt = strings.Replace(fullTextFmt, "\u003cipv6\u003e", ip.String()[0:3], -1)
				fullTextFmt = strings.Replace(fullTextFmt, "\u003ccidr6\u003e", addr.String(), -1)
			}
		}
	}

	// setting up/down flag
	fullTextFmt = strings.Replace(fullTextFmt, "\u003cstatus\u003e", status, -1)

	// clearing unset fields i.e. because of ipv6 single-stack
	fullTextFmt = strings.Replace(fullTextFmt, "\u003cipv4\u003e", "", -1)
	fullTextFmt = strings.Replace(fullTextFmt, "\u003ccidr4\u003e", "", -1)
	fullTextFmt = strings.Replace(fullTextFmt, "\u003cipv6\u003e", "", -1)
	fullTextFmt = strings.Replace(fullTextFmt, "\u003ccidr6\u003e", "", -1)
	fullTextFmt = strings.Replace(fullTextFmt, "\u003clocal6\u003e", "", -1)

	// removing the last %s placeholder from final string
	b.FullText = fmt.Sprintf(fullTextFmt, "")
}
