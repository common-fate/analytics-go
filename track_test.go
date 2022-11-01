package analytics

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/common-fate/analytics-go/acore"
	"github.com/stretchr/testify/assert"
)

func TestRequestCreated(t *testing.T) {
	tests := []struct {
		name       string
		ref        string
		data       Event
		deployment *Deployment
	}{
		{
			name: "ok",
			ref:  strings.TrimSpace(fixture("request-created.json")),
			data: &RequestCreated{
				RequestedBy: "usr_123",
				Provider:    "commonfate/test-provider@v1",
				Rule:        "rul_123",
				Timing: Timing{
					Mode:            TimingModeASAP,
					DurationSeconds: 100,
				},
				HasReason: true,
			},
			deployment: &Deployment{
				ID:         "dep_123",
				Version:    "v0.9.0",
				UserCount:  100,
				GroupCount: 10,
			},
		},
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

			ReplaceGlobal(client)
			SetDeploymentLoader(&testLoader{Deployment: tt.deployment})
			Track(context.Background(), tt.data)
			Close()

			res := string(<-body)
			assert.Equal(t, tt.ref, res)
		})
	}
}

func fixture(name string) string {
	f, err := os.Open(filepath.Join("fixtures", name))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return string(b)
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
