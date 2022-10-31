package acore

import "runtime/debug"

// Version of the client.
func getLibraryVersion() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}

	for _, dep := range bi.Deps {
		if dep.Path == "github.com/common-fate/analytics" {
			return dep.Version
		}
	}
	return "dev"
}
