package analytics

import (
	"os"
	"sync"

	"github.com/common-fate/analytics-go/acore"
	"go.uber.org/zap"
)

var (
	depMutex     = &sync.RWMutex{}
	globalDeploy *Deployment
)

// Deployment is a Common Fate deployment identifier.
// If you're editing this make sure you edit the Traits()
// method to ensure the properties propagate.
type Deployment struct {
	ID         string `json:"id"`
	Version    string `json:"version"`
	UserCount  int    `json:"user_count"`
	GroupCount int    `json:"group_count"`
}

// Traits returns the traits to use for the group identifier
func (d Deployment) Traits() acore.Traits {
	return acore.NewTraits().
		Set("version", d.Version).
		Set("user_count", d.UserCount).
		Set("group_count", d.GroupCount).
		Set("groupType", "deployment").
		Set("id", d.ID)
}

func getDeployment() *Deployment {
	depMutex.RLock()
	defer depMutex.RUnlock()
	d := globalDeploy
	return d
}

// SetDeployment sets deployment information.
func SetDeployment(dep *Deployment) {
	depMutex.Lock()
	defer depMutex.Unlock()
	globalDeploy = dep

	if os.Getenv("CF_ANALYTICS_DEBUG") == "true" {
		zap.L().Named("cf-analytics").Info("set deployment", zap.Any("deployment", dep))
	}
}
