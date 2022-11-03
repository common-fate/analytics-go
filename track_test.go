package analytics

import (
	"sync"
	"testing"

	"github.com/common-fate/analytics-go/acore"
	"github.com/stretchr/testify/assert"
)

// testEvent implements the Event interface for tests.
type testEvent struct {
	ExampleUserID string     `json:"example_user_id" analytics:"usr"`
	TestStruct    testStruct `json:"test_struct"`
}

type testStruct struct {
	Value string
}

func (e *testEvent) userID() string      { return e.ExampleUserID }
func (e *testEvent) Type() string        { return "TEST_ONLY" }
func (e *testEvent) EmittedWhen() string { return "during tests" }
func (e *testEvent) fixture()            {}

func TestClient_marshalToCapture(t *testing.T) {
	type fields struct {
		mu           *sync.Mutex
		deploymentID string
		coreclient   acore.Client
	}
	type args struct {
		e Event
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    acore.Capture
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				deploymentID: "dep_123",
			},
			args: args{
				e: &testEvent{
					ExampleUserID: "usr_123",
				},
			},
			want: acore.Capture{
				// DistinctId should always match the same user ID in properties,
				// otherwise there is a problem with the client-side hashing.
				DistinctId: "usr_-CHh8_rdIqAotcBsP64GKQkfzW2hb1JDJ_6u7q4zom4",
				Event:      "TEST_ONLY",
				Properties: acore.Properties{
					"example_user_id": "usr_-CHh8_rdIqAotcBsP64GKQkfzW2hb1JDJ_6u7q4zom4",
					"test_struct":     testStruct{},
				},
				Groups: acore.NewGroups().Set("deployment", "dep_123"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				mu:           tt.fields.mu,
				deploymentID: &tt.fields.deploymentID,
				coreclient:   tt.fields.coreclient,
			}
			got, err := c.marshalToCapture(tt.args.e)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.marshalToCapture() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
