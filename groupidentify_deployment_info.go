package analytics

import (
	"github.com/common-fate/analytics-go/acore"
)

func init() {
	registerEvent(&DeploymentInfo{})
}

type DeploymentInfo struct {
	ID         string `json:"-"`
	Version    string `json:"version"`
	UserCount  int    `json:"user_count"`
	GroupCount int    `json:"group_count"`
	IDP        string `json:"idp"`
	Stage      string `json:"stage,omitempty"` // dev, prod, uat, etc.
}

func (d *DeploymentInfo) userID() string { return "" }

func (d *DeploymentInfo) Type() string { return "cf:groupidentify:deployment" }

func (d *DeploymentInfo) EmittedWhen() string { return "Deployment updated" }

func (d *DeploymentInfo) marshalEvent(ctx marshalContext) ([]acore.Message, error) {
	m := []acore.Message{acore.GroupIdentify{
		Type:       "deployment",
		Key:        d.ID,
		Properties: eventToProperties(d),
	}}

	return m, nil
}

func (d *DeploymentInfo) fixture() {
	*d = DeploymentInfo{
		ID:         "dep_123",
		Version:    "v0.9.0",
		UserCount:  10,
		GroupCount: 5,
		IDP:        "cognito",
		Stage:      "test",
	}
}
