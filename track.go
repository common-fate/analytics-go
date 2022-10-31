package analytics

import (
	"context"

	"github.com/common-fate/analytics-go/acore"
)

// Event is a product analytics event that is tracked.
type Event interface {
	userID() string
	eventType() string
}

// Track an event using the global analytics client.
func Track(ctx context.Context, e Event) {
	e = hashValues(e)
	go func() {

		uid := e.userID()
		typ := e.eventType()
		client := G()
		dep := globalDeployment.Get(ctx)

		evt := acore.Track{
			UserId:     uid,
			Event:      typ,
			Properties: e,
		}
		if dep != nil {
			client.Enqueue(acore.Group{
				GroupId: dep.ID,
				Traits:  dep.Traits(),
				UserId:  uid,
			})
		}

		client.Enqueue(evt)
	}()
}
