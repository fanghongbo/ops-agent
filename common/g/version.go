package g

import "fmt"

var (
	Version    string = "v0.1"
	BinaryName string = "ops-agent.dev"
)

func VersionInfo() string {
	return fmt.Sprintf("%s", Version)
}

func AgentInfo() string {
	return fmt.Sprintf("%s.%s", BinaryName, Version)
}
