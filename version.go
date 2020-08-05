package main

import "github.com/fanghongbo/ops-agent/common/g"

var (
	Version    = "v1.0"
	BinaryName = "ops-agent"
)

func init() {
	g.BinaryName = BinaryName
	g.Version = Version
}
