package main

import (
	_ "github.com/coredns/coredns/core/plugin"
	_ "github.com/iandri/cname"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/coremain"
)

func init() {
	// Insert cname directive before file directive
	var i int
	for i = 0; i < len(dnsserver.Directives); i++ {
		if dnsserver.Directives[i] == "file" {
			break
		}
	}
	dnsserver.Directives = append(dnsserver.Directives, "")
	copy(dnsserver.Directives[i+1:], dnsserver.Directives[i:])
	dnsserver.Directives[i] = "cname"
}

func main() {
	coremain.Run()
}
