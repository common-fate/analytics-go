package main

import (
	"github.com/common-fate/analytics-go"
	"go.uber.org/zap"
)

func main() {
	analytics.Configure(analytics.Development)
	defer analytics.Close()

	log := zap.Must(zap.NewDevelopment())
	zap.ReplaceGlobals(log)

	analytics.SetDeployment(&analytics.Deployment{
		ID:         "dep_123",
		Version:    "v0.0.0",
		UserCount:  10,
		GroupCount: 10,
	})

	analytics.Track(&analytics.RequestCreated{
		RequestedBy: "usr_123",
		Provider:    "commonfate/test-provider@v1",
		Rule:        "rul_123",
		Timing: analytics.Timing{
			DurationSeconds: 100,
			Mode:            analytics.TimingModeASAP,
		},
		HasReason: true,
	})
}
