package analytics

func init() {
	registerEvent(&RuleCreated{})
}

type RuleCreated struct {
	RuleID                string `json:"rule_id" analytics:"rul"`
	CreatedBy             string `json:"created_by" analytics:"usr"`
	UsesSelectableOptions bool   `json:"uses_selectable_options"`
	UsesDynamicOptions    bool   `json:"uses_dynamic_options"`
	Provider              string `json:"provider"`
	MaxDurationSeconds    int    `json:"max_duration_seconds"`
	RequiresApproval      bool   `json:"requires_approval"`
}

func (r *RuleCreated) userID() string { return r.CreatedBy }

func (r *RuleCreated) Type() string { return "cf:rule.created" }

func (r *RuleCreated) EmittedWhen() string { return "Access Rule was created" }

func (r *RuleCreated) fixture() {
	*r = RuleCreated{
		RuleID:                "rul_123",
		CreatedBy:             "usr_123",
		UsesSelectableOptions: true,
		UsesDynamicOptions:    true,
		Provider:              "commonfate/test-provider@v1",
		MaxDurationSeconds:    100,
		RequiresApproval:      true,
	}
}
