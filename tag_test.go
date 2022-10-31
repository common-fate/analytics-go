package analytics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testTaggedStruct struct {
	UserID string `analytics:"usr"`
}

func Test_hashValues(t *testing.T) {
	tests := []struct {
		name  string
		event any
		want  any
	}{
		{
			name: "ok",
			event: &testTaggedStruct{
				UserID: "something",
			},
			want: &testTaggedStruct{
				UserID: "usr_kh6XMrBb8gxpzREC4mAOSNs862lIy8tjE9fNDBWrjRE",
			},
		},
		{
			name: "hash of empty value should still be empty",
			event: &testTaggedStruct{
				UserID: "",
			},
			want: &testTaggedStruct{
				UserID: "",
			},
		},
		{
			name: "request created event",
			event: &RequestCreated{
				RequestedBy: "usr_123",
			},
			want: &RequestCreated{
				RequestedBy: "usr_-CHh8_rdIqAotcBsP64GKQkfzW2hb1JDJ_6u7q4zom4",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hashValues(tt.event)
			assert.Equal(t, tt.want, got)
		})
	}
}
