package analytics

func init() {
	registerEvent(&RequestReviewed{})
}

type RequestReviewed struct {
	RequestedBy    string  `json:"requested_by" analytics:"usr"`
	ReviewedBy     string  `json:"reviewed_by" analytics:"usr"`
	AccessGroupID  string  `json:"access_group_id"`
	RequestID      string  `json:"request_id"`
	TargetsCount   int     `json:"targets_count"`
	Timing         Timing  `json:"timing"`
	OverrideTiming *Timing `json:"override_timing"`
	HasReason      bool    `json:"has_reason"`
	// PendingDurationSeconds is how long the request has been waiting for a review.
	PendingDurationSeconds float64 `json:"pending_duration_seconds"`
	// Review is APPROVE or DENY
	Review          string `json:"review"`
	ReviewerIsAdmin bool   `json:"reviewer_is_admin"`
}

// func (r *RequestReviewed) payloads() []acore.Message {
// 	return []acore.Message{acore.Identify{
// 		UserId: r.ReviewedBy,
// 		Traits: acore.NewTraits().Set("role", getRole(r.ReviewerIsAdmin)),
// 	}}
// }

func (r *RequestReviewed) userID() string { return r.ReviewedBy }

func (r *RequestReviewed) Type() string { return "cf:request.reviewed" }

func (r *RequestReviewed) EmittedWhen() string { return "Access Request was reviewed" }

func (r *RequestReviewed) fixture() {
	*r = RequestReviewed{
		RequestedBy:   "usr_123",
		ReviewedBy:    "usr_234",
		RequestID:     "req_123",
		AccessGroupID: "group_123",
		TargetsCount:  4,
		OverrideTiming: &Timing{
			Mode:            TimingModeScheduled,
			DurationSeconds: 50,
		},
		PendingDurationSeconds: 200,
		Review:                 "APPROVE",
		ReviewerIsAdmin:        true,
		Timing: Timing{
			Mode:            TimingModeASAP,
			DurationSeconds: 100,
		},
		HasReason: true,
	}
}
