package analytics

import (
	"os"

	"github.com/common-fate/analytics-go/acore"
	"go.uber.org/zap"
)

// Event is a product analytics event that is tracked.
type Event interface {
	userID() string
	eventType() string
}

// Track an event using the global analytics client.
func Track(e Event) {
	e = hashValues(e)
	go func() {
		uid := e.userID()
		typ := e.eventType()
		client := G()
		dep := getDeployment()

		evt := acore.Track{
			UserId:     uid,
			Event:      typ,
			Properties: e,
		}
		if dep != nil {
			enqueueAndLog(client, acore.Group{
				GroupId: dep.ID,
				Traits:  dep.Traits(),
				UserId:  uid,
			})
		}

		enqueueAndLog(client, evt)
	}()
}

// enqueueAndLog logs the analytics event using the global zap logger if CF_ANALYTICS_DEBUG is set.
func enqueueAndLog(c acore.Client, m acore.Message) {
	if os.Getenv("CF_ANALYTICS_DEBUG") == "true" {
		zap.L().Named("cf-analytics").Info("emitting analytics event", zap.String("url", c.EndpointURL()), zap.Any("event", m))
	}

	c.Enqueue(m)
}
