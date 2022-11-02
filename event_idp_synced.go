package analytics

func init() {
	registerEvent(&IDPSynced{})
}

type IDPSynced struct {
	DeploymentID string `json:"deployment_id"`
	UserCount    int    `json:"user_count" analytics:"rul"`
	GroupCount   int    `json:"group_count" analytics:"usr"`
	IDP          string `json:"idp"`
}

func (e *IDPSynced) userID() string { return e.DeploymentID }

func (r *IDPSynced) Type() string { return "cf:idp.synced" }

func (r *IDPSynced) deploymentEvent() bool { return true }

func (r *IDPSynced) EmittedWhen() string { return "IDP was synced" }

func (r *IDPSynced) fixture() {
	*r = IDPSynced{
		UserCount:    10,
		GroupCount:   100,
		IDP:          "cognito",
		DeploymentID: "dep_123",
	}
}
