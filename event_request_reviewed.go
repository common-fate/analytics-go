package analytics

type RequestReviewed struct {
	RequestedBy    string  `json:"requested_by" analytics:"usr"`
	ReviewedBy     string  `json:"reviewed_by" analytics:"usr"`
	Provider       string  `json:"provider"`
	Rule           string  `json:"rule" analytics:"rul"`
	Timing         Timing  `json:"timing"`
	OverrideTiming *Timing `json:"override_timing"`
	HasReason      bool    `json:"has_reason"`
	// PendingDurationSeconds is how long the request has been waiting for a review.
	PendingDurationSeconds float64 `json:"pending_duration_seconds"`
	// Review is APPROVE or DENY
	Review string `json:"review"`
}

func (r *RequestReviewed) userID() string { return r.ReviewedBy }

func (r *RequestReviewed) eventType() string { return "cf:request.reviewed" }
