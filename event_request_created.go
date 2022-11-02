package analytics

func init() {
	registerEvent(&RequestCreated{})
}

type TimingMode string

var (
	TimingModeASAP      TimingMode = "asap"
	TimingModeScheduled TimingMode = "scheduled"
)

type Timing struct {
	Mode            TimingMode `json:"mode"`
	DurationSeconds float64    `json:"duration_seconds"`
}

type RequestCreated struct {
	RequestedBy string `json:"requested_by" analytics:"usr"`
	Provider    string `json:"provider"`
	RuleID      string `json:"rule_id" analytics:"rul"`
	Timing      Timing `json:"timing"`
	HasReason   bool   `json:"has_reason"`
}

func (r *RequestCreated) userID() string { return r.RequestedBy }

func (r *RequestCreated) Type() string { return "cf:request.created" }

func (r *RequestCreated) EmittedWhen() string { return "Access Request was created" }

func (r *RequestCreated) fixture() {
	*r = RequestCreated{
		RequestedBy: "usr_123",
		Provider:    "commonfate/test-provider@v1",
		RuleID:      "rul_123",
		Timing: Timing{
			Mode:            TimingModeASAP,
			DurationSeconds: 100,
		},
		HasReason: true,
	}
}
