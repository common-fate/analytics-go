package acore

import "runtime/debug"

// Version of the client.
func getVersion() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}

	for _, dep := range bi.Deps {
		if dep.Path == "github.com/common-fate/analytics-go" {
			return dep.Version
		}
	}
	return "dev"
}
