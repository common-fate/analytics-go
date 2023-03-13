package analytics

func init() {
	registerEvent(&RuleUpdated{})
}

type RuleUpdated struct {
	RuleID                string   `json:"rule_id" analytics:"rul"`
	UpdatedBy             string   `json:"updated_by" analytics:"usr"`
	UsesSelectableOptions bool     `json:"uses_selectable_options"`
	UsesDynamicOptions    bool     `json:"uses_dynamic_options"`
	Provider              Provider `json:"provider"`
	BuiltInProvider       string   `json:"built_in_provider"`
	PDKProvider           bool     `json:"pdk_provider"`
	MaxDurationSeconds    int      `json:"max_duration_seconds"`
	RequiresApproval      bool     `json:"requires_approval"`
	OptionCount           int      `json:"option_count"`
}

func (r *RuleUpdated) userID() string { return r.UpdatedBy }

func (r *RuleUpdated) Type() string { return "cf:rule.updated" }

func (r *RuleUpdated) EmittedWhen() string { return "Access Rule was updated" }

func (r *RuleUpdated) fixture() {
	*r = RuleUpdated{
		RuleID:                "rul_123",
		UpdatedBy:             "usr_123",
		UsesSelectableOptions: true,
		UsesDynamicOptions:    true,
		BuiltInProvider:       "commonfate/test-provider@v1",
		MaxDurationSeconds:    100,
		RequiresApproval:      true,
		OptionCount:           5,
	}
}
