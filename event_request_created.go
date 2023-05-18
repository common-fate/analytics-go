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
	RequestID   string `json:"request_id"`
	HasReason   bool   `json:"has_reason"`
	// total number of resources selected.
	TargetsCount int `json:"targets_count"`
	// total number of access groups selected.
	AccessGroupsCount int `json:"access_groups_count"`
}

func (r *RequestCreated) userID() string { return r.RequestedBy }

func (r *RequestCreated) Type() string { return "cf:request.created" }

func (r *RequestCreated) EmittedWhen() string { return "Access Request was created" }

func (r *RequestCreated) fixture() {
	*r = RequestCreated{
		RequestedBy:       "usr_123",
		TargetsCount:      6,
		AccessGroupsCount: 3,
		RequestID:         "req_123",
		Timing: Timing{
			Mode:            TimingModeASAP,
			DurationSeconds: 100,
		},
		HasReason: true,
	}
}
