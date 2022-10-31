package main

import (
	"context"

	"github.com/common-fate/analytics"
	"go.uber.org/zap"
)

func main() {
	analytics.SetDevelopment(true)
	defer analytics.Close()

	log := zap.Must(zap.NewDevelopment())
	zap.ReplaceGlobals(log)

	analytics.SetDeploymentLoader(func(ctx context.Context) (*analytics.Deployment, error) {
		d := analytics.Deployment{
			ID:         "dep_123",
			Version:    "v0.0.0",
			UserCount:  10,
			GroupCount: 10,
		}
		return &d, nil
	})

	analytics.Track(context.Background(), &analytics.RequestCreated{
		RequestedBy: "usr_123",
		Provider:    "commonfate/test-provider@v1",
		Rule:        "rul_123",
		Duration:    100,
		TimingMode:  analytics.TimingModeASAP,
		HasReason:   true,
	})
}
