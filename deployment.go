package analytics

import (
	"os"

	"github.com/common-fate/analytics-go/acore"
	"go.uber.org/zap"
)

// Deployment is a Common Fate deployment identifier.
// If you're editing this make sure you edit the Traits()
// method to ensure the properties propagate.
type Deployment struct {
	ID         string `json:"id"`
	Version    string `json:"version"`
	UserCount  int    `json:"user_count"`
	GroupCount int    `json:"group_count"`
	Stage      string `json:"stage"` // dev, prod, uat, etc.
}

// Traits returns the traits to use for the group identifier
func (d Deployment) Traits() acore.Traits {
	return acore.NewTraits().
		Set("version", d.Version).
		Set("user_count", d.UserCount).
		Set("group_count", d.GroupCount).
		Set("stage", d.Stage).
		Set("groupType", "deployment").
		Set("id", d.ID)
}

// SetDeployment sets deployment information.
func (c *Client) SetDeployment(dep *Deployment) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.deployment = dep

	if os.Getenv("CF_ANALYTICS_DEBUG") == "true" {
		zap.L().Named("cf-analytics").Info("set deployment", zap.Any("deployment", dep))
	}
}
