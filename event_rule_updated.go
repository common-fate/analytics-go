package analytics

func init() {
	registerEvent(&RuleUpdated{})
}

type RuleUpdated struct {
	RuleID              string   `json:"rule_id" analytics:"rul"`
	UpdatedBy           string   `json:"updated_by" analytics:"usr"`
	MaxDurationSeconds  int      `json:"max_duration_seconds"`
	RequiresApproval    bool     `json:"requires_approval"`
	HasFilterExpression bool     `json:"has_filter_expression"`
	TargetsCount        int      `json:"targets_count"`
	Targets             []string `json:"targets"`
}

func (r *RuleUpdated) userID() string { return r.UpdatedBy }

func (r *RuleUpdated) Type() string { return "cf:rule.updated" }

func (r *RuleUpdated) EmittedWhen() string { return "Access Rule was updated" }

func (r *RuleUpdated) fixture() {
	*r = RuleUpdated{
		RuleID:             "rul_123",
		UpdatedBy:          "usr_123",
		Targets:            []string{"commonfate/test-provider1@v2", "commonfate/test-provider2@v3"},
		TargetsCount:       2,
		MaxDurationSeconds: 100,
		RequiresApproval:   true,
	}
}
