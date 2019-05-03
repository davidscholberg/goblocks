package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"io/ioutil"
	"net/http"
	"strings"
)

// PublicIp represents the configuration for public ip display block.
type PublicIp struct {
	BlockConfigBase `yaml:",inline"`
	IpFormat        string `yaml:"ip_format"`
}

func (c PublicIp) UpdateBlock(b *i3barjson.Block) {
	b.Color = c.Color

	switch c.IpFormat {
	case "city":
		b.FullText = fmt.Sprintf("%s%s", c.Label, get("https://ifconfig.co/city"))
	default:
		b.FullText = fmt.Sprintf("%s%s", c.Label, get("https://ifconfig.co/ip"))
	}
}

func get(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		return err.Error()
	}

	ip, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return strings.TrimSpace(string(ip))
}
