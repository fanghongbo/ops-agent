package g

import "fmt"

var (
	BinaryName string = "ops-agent.dev"
	Version    string = "v0.1"
)

func VersionInfo() string {
	return fmt.Sprintf("%s", Version)
}

func AgentInfo() string {
	return fmt.Sprintf("%s.%s", BinaryName, Version)
}
