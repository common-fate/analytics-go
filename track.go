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
func (c *Client) Track(e Event) {
	e = hashValues(e)
	uid := e.userID()
	typ := e.eventType()

	evt := acore.Track{
		UserId:     uid,
		Event:      typ,
		Properties: e,
	}
	if c.deployment != nil {
		enqueueAndLog(c.coreclient, acore.Group{
			GroupId: c.deployment.ID,
			Traits:  c.deployment.Traits(),
			UserId:  uid,
		})
	}

	enqueueAndLog(c.coreclient, evt)
}

// enqueueAndLog logs the analytics event using the global zap logger if CF_ANALYTICS_DEBUG is set.
func enqueueAndLog(c acore.Client, m acore.Message) {
	if os.Getenv("CF_ANALYTICS_DEBUG") == "true" {
		zap.L().Named("cf-analytics").Info("emitting analytics event", zap.String("url", c.EndpointURL()), zap.Any("event", m))
	}

	c.Enqueue(m)
}
