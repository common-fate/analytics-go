package analytics

func init() {
	registerEvent(&IDPSynced{})
}

type IDPSynced struct {
	UserCount  int    `json:"user_count" analytics:"rul"`
	GroupCount int    `json:"group_count" analytics:"usr"`
	IDP        string `json:"idp"`
}

func (r *IDPSynced) userID() string { return "" }

func (r *IDPSynced) Type() string { return "cf:idp.synced" }

func (r *IDPSynced) EmittedWhen() string { return "IDP was synced" }

func (r *IDPSynced) fixture() {
	*r = IDPSynced{
		UserCount:  10,
		GroupCount: 100,
		IDP:        "cognito",
	}
}
