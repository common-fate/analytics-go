package analytics

import (
	"testing"

	"github.com/common-fate/analytics-go/acore"
	"github.com/stretchr/testify/assert"
)

func Test_eventToProperties(t *testing.T) {
	tests := []struct {
		name  string
		event Event
		want  acore.Properties
	}{
		{
			name: "ok",
			event: &RequestCreated{
				RequestedBy: "usr_123",
			},
			want: acore.NewProperties().
				Set("has_reason", false).
				Set("rule_id", "").
				Set("provider", "").
				Set("requested_by", "usr_-CHh8_rdIqAotcBsP64GKQkfzW2hb1JDJ_6u7q4zom4").
				Set("timing", Timing{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := eventToProperties(tt.event)
			assert.Equal(t, tt.want, got)
		})
	}
}
