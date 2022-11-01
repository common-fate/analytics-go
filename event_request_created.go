package analytics

type TimingMode string

var (
	TimingModeASAP      TimingMode = "asap"
	TimingModeScheduled TimingMode = "scheduled"
)

type RequestCreated struct {
	RequestedBy     string     `json:"requested_by" analytics:"usr"`
	Provider        string     `json:"provider"`
	Rule            string     `json:"rule" analytics:"rul"`
	DurationSeconds float64    `json:"duration_seconds"`
	TimingMode      TimingMode `json:"timing_mode"`
	HasReason       bool       `json:"has_reason"`
}

func (r *RequestCreated) userID() string { return r.RequestedBy }

func (r *RequestCreated) eventType() string { return "cf:request.created" }
