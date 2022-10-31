# analytics

This repository contains product analytics as defined in [RFD#8](https://github.com/common-fate/rfds/discussions/8).

## Event types

The [fixtures](./fixtures/) folder contains examples of each event type which is dispatched, along with the properties.

## Usage

The `acore` package contains the core analytics client. The client is forked from the Rudderstack Go client (which itself appears to have been forked from Segment's Go SDK).

```go
package main

import (
    "github.com/rudderlabs/analytics-go"
)

func main() {
    // set development or production mode for the analytics client.
    analytics.SetDevelopment(true)

    // analytics.Close() must be called prior to exiting to flush any pending messages.
	defer analytics.Close()

    // Set a loader function for the deployment.
	analytics.SetDeploymentLoader(func(ctx context.Context) (*analytics.Deployment, error) {
		d := analytics.Deployment{
			ID:         "dep_123",
			Version:    "v0.0.0",
			UserCount:  10,
			GroupCount: 10,
		}
		return &d, nil
	})

    // Track an event.
	analytics.Track(context.Background(), analytics.RequestCreated{
		RequestedBy: "usr_123",
		Provider:    "commonfate/test-provider@v1",
		Rule:        "rul_123",
		Duration:    100,
		TimingMode:  analytics.TimingModeASAP,
		HasReason:   true,
	})
}
```

The library handles client-side hashing identifiers such as `usr_123`. We transform `usr_123` using a SHA256 hash into `usr_-CHh8_rdIqAotcBsP64GKQkfzW2hb1JDJ_6u7q4zom4` prior to events being dispatched. In the library this is controlled by the `analytics:"usr"` struct tag added to the event.
