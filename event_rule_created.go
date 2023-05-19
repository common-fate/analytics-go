package analytics

func init() {
	registerEvent(&RuleCreated{})
}

type RuleCreated struct {
	RuleID              string   `json:"rule_id" analytics:"rul"`
	CreatedBy           string   `json:"created_by" analytics:"usr"`
	MaxDurationSeconds  int      `json:"max_duration_seconds"`
	RequiresApproval    bool     `json:"requires_approval"`
	HasFilterExpression bool     `json:"has_filter_expression"`
	TargetsCount        int      `json:"targets_count"`
	Targets             []string `json:"targets"`
}

func (r *RuleCreated) userID() string { return r.CreatedBy }

func (r *RuleCreated) Type() string { return "cf:rule.created" }

func (r *RuleCreated) EmittedWhen() string { return "Access Rule was created" }

func (r *RuleCreated) fixture() {
	*r = RuleCreated{
		RuleID:              "rul_123",
		CreatedBy:           "usr_123",
		Targets:             []string{"commonfate/test-provider@v1", "commonfate/test-provider2@v2"},
		HasFilterExpression: true,
		TargetsCount:        2,
		MaxDurationSeconds:  100,
		RequiresApproval:    true,
	}
}
