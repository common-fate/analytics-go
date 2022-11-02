package main

import (
	"github.com/common-fate/analytics-go"
	"go.uber.org/zap"
)

func main() {
	c := analytics.New(analytics.Development)
	defer c.Close()

	log := zap.Must(zap.NewDevelopment())
	zap.ReplaceGlobals(log)

	c.SetDeployment(&analytics.Deployment{
		ID:      "dep_100",
		Version: "v0.0.0",
	})

	c.Track(&analytics.RequestCreated{
		RequestedBy: "usr_500",
		Provider:    "commonfate/test-provider@v1",
		RuleID:      "rul_123",
		Timing: analytics.Timing{
			DurationSeconds: 100,
			Mode:            analytics.TimingModeASAP,
		},
		HasReason: true,
	})

	c.Track(&analytics.IDPSynced{
		UserCount:  0,
		GroupCount: 0,
		IDP:        "cognito",
	})
}
