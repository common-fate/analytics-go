package analytics

func init() {
	registerEvent(&RuleArchived{})
}

type RuleArchived struct {
	RuleID     string `json:"rule_id" analytics:"rul"`
	ArchivedBy string `json:"archived_by" analytics:"usr"`
}

func (r *RuleArchived) userID() string { return r.ArchivedBy }

func (r *RuleArchived) Type() string { return "cf:rule.archived" }

func (r *RuleArchived) EmittedWhen() string { return "Access Rule was archived" }

func (r *RuleArchived) fixture() {
	*r = RuleArchived{
		RuleID:     "rul_123",
		ArchivedBy: "usr_123",
	}
}
