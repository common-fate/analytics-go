package analytics

func init() {
	registerEvent(&RequestRevoked{})
}

type RequestRevoked struct {
	RequestedBy string `json:"requested_by" analytics:"usr"`
	RevokedBy   string `json:"revoked_by" analytics:"usr"`
	AccessGroup string `json:"access_group"`
	Timing      Timing `json:"timing"`
	HasReason   bool   `json:"has_reason"`
}

func (r *RequestRevoked) userID() string { return r.RevokedBy }

func (r *RequestRevoked) Type() string { return "cf:request.revoked" }

func (r *RequestRevoked) EmittedWhen() string { return "Access Request was revoked" }

func (r *RequestRevoked) fixture() {
	*r = RequestRevoked{
		RequestedBy: "usr_123",
		RevokedBy:   "usr_234",
		AccessGroup: "rul_123",
		Timing: Timing{
			Mode:            TimingModeASAP,
			DurationSeconds: 100,
		},
		HasReason: true,
	}
}
