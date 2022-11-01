package analytics

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
	Rule        string `json:"rule" analytics:"rul"`
	Timing      Timing `json:"timing"`
	HasReason   bool   `json:"has_reason"`
}

func (r *RequestCreated) userID() string { return r.RequestedBy }

func (r *RequestCreated) eventType() string { return "cf:request.created" }
