package analytics

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/common-fate/analytics-go/acore"
	"github.com/common-fate/analytics-go/internal"
	"github.com/stretchr/testify/assert"
)

// if -update is provided, the fixture files will be updated.
var update = flag.Bool("update", false, "update fixture files")

// TestEventsSuite checks that a fixture exists for each event,
// and that the emitted event matches the fixture.
//
// To update/generate fixture data, run this test with the -update
// flag:
//
//	go test -update
func TestEventsSuite(t *testing.T) {
	type testcase struct {
		name        string
		fixturepath string
		want        string
		event       Event
		deployment  *Deployment
	}

	var tests []testcase

	// create a test case for each event type.
	for _, e := range AllEvents {
		// replace cf:request.created with cf-request-created
		f := internal.FixturePath(e.Type())

		e.fixture()

		tc := testcase{
			name:        e.Type(),
			fixturepath: f,
			event:       e,
			// always use this as an example deployment.
			deployment: &Deployment{
				ID:         "dep_123",
				Version:    "v0.9.0",
				UserCount:  100,
				GroupCount: 10,
			},
		}

		if !*update {
			// If we're updating the fixture, don't try and load the existing one.
			// It might not exist yet.
			tc.want = loadFixture(t, f)
		}

		tests = append(tests, tc)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, server := mockServer()
			defer server.Close()

			client := acore.NewTestWithConfig(t, acore.Config{
				Endpoint:  server.URL,
				Verbose:   true,
				Interval:  time.Second * 10,
				Logger:    t,
				BatchSize: 10,
			}, acore.TestConfig{
				Now: mockTime,
				UID: mockId,
			})

			c := newClient(client)
			c.SetDeployment(tt.deployment)
			c.Track(tt.event)
			c.Close()

			res := string(<-body)

			if *update {
				res = res + "\n" // add newline to avoid code editors autoformatting one in.
				err := os.WriteFile(tt.fixturepath, []byte(res), 0644)
				if err != nil {
					t.Fatal(err)
				}
			} else {
				assert.Equal(t, tt.want, res)
			}
		})
	}
}

func loadFixture(t *testing.T, fixturepath string) string {
	f, err := os.Open(fixturepath)
	if err != nil {
		t.Fatalf("error loading fixture file (make sure you've created it!): %s", err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("error parsing fixture file: %s", err)
	}
	return strings.TrimSpace(string(b))
}

func mockServer() (chan []byte, *httptest.Server) {
	done := make(chan []byte, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := bytes.NewBuffer(nil)
		io.Copy(buf, r.Body)

		var v interface{}
		err := json.Unmarshal(buf.Bytes(), &v)
		if err != nil {
			panic(err)
		}

		b, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			panic(err)
		}

		done <- b
	}))

	return done, server
}

func mockId() string { return "KSUID" }

func mockTime() time.Time {
	return time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
}
