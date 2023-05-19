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

	c.SetDeploymentID("dep_101")

	c.Track(&analytics.RequestCreated{
		RequestID:         "req_123",
		RequestedBy:       "usr_501",
		TargetsCount:      4,
		AccessGroupsCount: 2,
		HasReason:         true,
	})

	c.Track(&analytics.UserInfo{
		ID:         "usr_501",
		IsAdmin:    false,
		GroupCount: 5,
	})

	c.Track(&analytics.DeploymentInfo{
		ID:         "dep_101",
		Version:    "v0.0.0",
		UserCount:  10,
		GroupCount: 1,
		IDP:        "cognito",
		Stage:      "dev",
	})
}
