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
			event: &testEvent{
				ExampleUserID: "usr_123",
				TestStruct: testStruct{
					Value: "123",
				},
			},
			want: acore.NewProperties().
				Set("example_user_id", "usr_-CHh8_rdIqAotcBsP64GKQkfzW2hb1JDJ_6u7q4zom4").
				Set("test_struct", testStruct{
					Value: "123",
				}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := eventToProperties(tt.event)
			assert.Equal(t, tt.want, got)
		})
	}
}
